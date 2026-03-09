package crawler

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/Dercraker/SearchEngine/internal/DAL"
	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
	"github.com/Dercraker/SearchEngine/internal/crawler/storage"
	"github.com/Dercraker/SearchEngine/internal/shared/instanceId"
	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
	"github.com/google/uuid"
)

type QueueRunner struct {
	Logger    *slog.Logger
	Processor UrlProcessor
	Queue     storage.QueueStore
	Stats     *obs.Stats

	CanonicalOptions seeds.CanonicalOptions

	batchSize      int32
	StaleAfter     time.Duration
	MaxPagesPerRun int64

	InstanceId string
	Worker     int
}

type queueJob struct {
	Item DAL.ClaimNextBatchRow
	url  *url.URL
}

type queueResult struct {
	Item DAL.ClaimNextBatchRow
	err  error
}

func (r *QueueRunner) RunOnce(ctx context.Context) (*obs.Stats, error) {
	rid := uuid.NewString()
	ctx = requestId.WithRunId(ctx, rid)
	ctx = instanceId.WithInstanceId(ctx, r.InstanceId)

	r.Stats.StartTime = time.Now()
	r.Logger.Info(
		string(obs.RunStart),
		slog.String("request_id", rid),
		slog.String("instance_id", r.InstanceId),
		slog.Int("workers", r.Worker),
		slog.Duration("stale_after", r.StaleAfter),
	)

	if r.StaleAfter > 0 {
		if err := r.Queue.ReleaseStale(ctx, r.StaleAfter); err != nil {
			r.Stats.DBFailed.Add(1)
		}
	}

	for {
		if r.MaxPagesPerRun > 0 && r.Stats.Processed.Load() >= r.MaxPagesPerRun {
			break
		}

		batch, qerr := r.Queue.ClaimNextBatch(ctx, r.batchSize)
		if qerr != nil {
			r.Stats.EndTime = time.Now()
			return r.Stats, qerr
		}
		if len(batch) == 0 {
			break
		}

		r.Logger.Info(string(obs.QueueClaim),
			slog.String("request_id", rid),
			slog.String("instance_id", r.InstanceId),
			slog.Int("claimed", len(batch)),
			slog.Int("batch_size", int(r.batchSize)),
		)

		r.processBatch(ctx, batch)
	}

	r.Stats.EndTime = time.Now()
	r.Logger.Info(string(obs.RunEnd),
		slog.String("request_id", rid),
		slog.Float64("duration_ms", r.Stats.DurationMs()),
		slog.Float64("duration_s", r.Stats.DurationS()),
		slog.Float64("duration_m", r.Stats.DurationM()),
		slog.Int64("processed", r.Stats.Processed.Load()),
		slog.Int64("success", r.Stats.Success.Load()),
		slog.Int64("failed", r.Stats.Failed.Load()),
		slog.Int64("inserted", r.Stats.Inserted.Load()),
		slog.Int64("updated", r.Stats.Updated.Load()),
		slog.Int64("unchanged", r.Stats.Unchanged.Load()),
		slog.Int64("touched", r.Stats.Touched.Load()),
		slog.Int64("skipped_non_html", r.Stats.SkippedNonHTML.Load()),
		slog.Int64("fetch_failed", r.Stats.FetchFailed.Load()),
		slog.Int64("db_failed", r.Stats.DBFailed.Load()),
		slog.Int64("retries", r.Stats.Retries.Load()),
	)
	return r.Stats, nil
}

func (r *QueueRunner) startWorkers(ctx context.Context, jobs <-chan queueJob, results chan<- queueResult) {
	workers := r.Worker
	batchSize := r.batchSize
	if workers <= 0 {
		workers = 1
	}
	if batchSize <= 0 {
		batchSize = int32(workers)
	}

	for i := 0; i < workers; i++ {
		go func(workerId int) {
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}

					err := r.Processor.Process(ctx, job.url)
					results <- queueResult{
						Item: job.Item,
						err:  err,
					}
				}
			}
		}(i)
	}
}

func buildJobs(batch []DAL.ClaimNextBatchRow) []queueJob {
	jobs := make([]queueJob, 0, len(batch))
	for i, item := range batch {
		pu, err := url.Parse(item.Url)
		if err != nil {
			continue
		}
		jobs = append(jobs, queueJob{
			Item: batch[i],
			url:  pu,
		})
	}
	return jobs
}

func (r *QueueRunner) processBatch(ctx context.Context, batch []DAL.ClaimNextBatchRow) {
	if len(batch) == 0 {
		return
	}

	jobs := make(chan queueJob, len(batch))
	results := make(chan queueResult, len(batch))

	r.startWorkers(ctx, jobs, results)

	preparedJobs := buildJobs(batch)

	skipped := len(batch) - len(preparedJobs)
	if skipped > 0 {
		r.Stats.SkippedNonHTML.Add(int64(skipped))
	}

	for _, job := range preparedJobs {
		select {
		case <-ctx.Done():
			return
		case jobs <- job:
		}
	}
	close(jobs)

	for i := 0; i < len(preparedJobs); i++ {
		select {
		case <-ctx.Done():
			return
		case res := <-results:
			r.handleResult(ctx, res)
		}
	}
}

func (r *QueueRunner) handleResult(ctx context.Context, res queueResult) {
	if res.err != nil {
		r.Stats.Failed.Add(1)

		next := nextRunAt(time.Now(), res.Item.Attempts)
		if err := r.Queue.MarkFailed(ctx, res.Item.Url, classifyLastError(res.err), next); err != nil {
			r.Stats.DBFailed.Add(1)
		}
		return
	}

	r.Stats.Success.Add(1)

	if err := r.Queue.MarkCrawled(ctx, res.Item.Url); err != nil {
		r.Stats.DBFailed.Add(1)
		return
	}

}

func nextRunAt(now time.Time, attempt int32) time.Time {
	var d time.Duration
	switch attempt {
	case 0:
		d = 30 * time.Second
	case 1:
		d = 1 * time.Minute
	case 2:
		d = 5 * time.Minute
	case 3:
		d = 15 * time.Minute
	case 4:
		d = 30 * time.Minute
	default:
		d = 1 * time.Hour
	}
	return now.Add(d)
}

func classifyLastError(err error) string {
	return err.Error()
}
