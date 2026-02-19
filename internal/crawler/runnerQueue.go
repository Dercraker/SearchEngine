package crawler

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
	"github.com/Dercraker/SearchEngine/internal/crawler/storage"
	"github.com/Dercraker/SearchEngine/internal/shared/customErrors"
	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
	"github.com/google/uuid"
)

type QueueRunner struct {
	Logger     *slog.Logger
	SeedSource seeds.Source
	Processor  UrlProcessor
	Queue      storage.QueueStore
	Stats      *obs.Stats

	CanonicalOptions seeds.CanonicalOptions

	batchSize      int32
	StaleAfter     time.Duration
	MaxPagesPerRun int64
}

func (r *QueueRunner) RunOnce(ctx context.Context) (*obs.Stats, error) {
	rid := uuid.NewString()
	ctx = requestId.WithRunId(ctx, rid)

	r.Stats.StartTime = time.Now()
	r.Logger.Info(
		string(obs.RunStart),
		slog.String("request_id", rid),
	)

	if r.StaleAfter > 0 {
		if err := r.Queue.ReleaseStale(ctx, r.StaleAfter); err != nil {
			r.Stats.DBFailed.Add(1)
		}
	}

	raw, err := r.SeedSource.Load(ctx)
	if err != nil {
		r.Stats.EndTime = time.Now()
		return r.Stats, err
	}
	list := seeds.SplitSeeds(raw)
	r.Stats.TotalSeeds = len(list)

	if len(list) == 0 {
		r.Logger.Error("[Crawler] no seeds provided")
		r.Stats.EndTime = time.Now()
		return r.Stats, err
	}

	for _, s := range list {
		u, nerr := seeds.NormalizeHTTPURL(s)
		if nerr != nil {
			r.Stats.InvalidSeeds++
			continue
		}

		key, kerr := seeds.CanonicalKey(u, r.CanonicalOptions)
		if kerr != nil {
			r.Stats.InvalidSeeds++
			continue
		}

		if err := r.Queue.Enqueue(ctx, key); err != nil {
			r.Stats.DBFailed.Add(1)
		}
	}

	for {
		if r.MaxPagesPerRun > 0 && r.Stats.Processed.Load() >= r.MaxPagesPerRun {
			break
		}

		items, qerr := r.Queue.ClaimNextBatch(ctx, r.batchSize)
		if qerr != nil {
			r.Stats.EndTime = time.Now()
			return r.Stats, qerr
		}
		if len(items) == 0 {
			break
		}

		r.Logger.Info(string(obs.QueueClaim),
			slog.String("request_id", rid),
			slog.Int("claimed", len(items)),
			slog.Int("batch_size", int(r.batchSize)),
		)

		for _, item := range items {
			urlStr := item.Url
			attemps := item.Attempts

			r.Stats.Processed.Add(1)

			pu, _ := url.Parse(urlStr)
			perr := r.Processor.Process(ctx, pu)
			if perr != nil {
				r.Stats.Failed.Add(1)

				next := nextRunAt(time.Now(), attemps)
				if err := r.Queue.MarkFailed(ctx, urlStr, classifyLastError(perr), next); err != nil {
					r.Stats.DBFailed.Add(1)
				}

				if errors.Is(perr, customErrors.ErrMaxPagesReached) {
					break
				}
				continue
			}
			r.Stats.Success.Add(1)
			if err := r.Queue.MarkCrawled(ctx, urlStr); err != nil {
				r.Stats.DBFailed.Add(1)
			}
		}
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
