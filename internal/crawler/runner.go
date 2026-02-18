package crawler

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
	"github.com/google/uuid"
)

type Runner struct {
	Logger     *slog.Logger
	SeedSource seeds.Source
	Processor  UrlProcessor
	Stats      *obs.Stats

	CanonicalOptions seeds.CanonicalOptions
}

func (r *Runner) RunOnce(ctx context.Context) (*obs.Stats, error) {
	r.Stats.StartTime = time.Now()

	rid := uuid.NewString()
	ctx = requestId.WithRunId(ctx, rid)

	raw, err := r.SeedSource.Load(ctx)
	if err != nil {
		r.Stats.EndTime = time.Now()
		return r.Stats, err
	}

	list := seeds.SplitSeeds(raw)

	r.Logger.Info(
		string(obs.RunStart),
		slog.String("request_id", rid),
		slog.Int("seeds_count", len(list)),
	)
	r.Stats.TotalSeeds = len(list)

	if len(list) == 0 {
		r.Logger.Error("[Crawler] no seeds provided")
		r.Stats.EndTime = time.Now()
		return r.Stats, err
	}

	seen := make(map[string]struct{}, len(list))

	for _, s := range list {
		u, nerr := seeds.NormalizeHTTPURL(s)
		if nerr != nil {
			r.Stats.InvalidSeeds++
			r.Logger.Error("[Crawler] invalid seed skipped", slog.Any("seed", s), slog.Any("error", nerr))
			continue
		}

		key, kerr := seeds.CanonicalKey(u, r.CanonicalOptions)
		if kerr != nil {
			r.Stats.InvalidSeeds++
			r.Logger.Error("[Crawler] invalid seed skipped", slog.Any("seed", s), slog.Any("error", kerr))
			continue
		}
		if _, ok := seen[key]; ok {
			r.Stats.DedupSkipped++
			continue
		}

		seen[key] = struct{}{}

		r.Logger.Info(
			string(obs.URLStart),
			slog.String("request_id", rid),
			slog.String("seed", key),
			slog.String("url", u.String()),
			slog.String("canonical_key", key),
		)
		cu, _ := url.Parse(key)

		r.Stats.Processed.Add(1)

		perr := r.Processor.Process(ctx, cu)
		if perr != nil {
			r.Stats.Failed.Add(1)
			continue
		}
		r.Stats.Success.Add(1)
	}
	r.Stats.EndTime = time.Now()
	r.Logger.Info(string(obs.RunEnd),
		slog.String("request_id", rid),
		slog.Float64("duration_ms", r.Stats.DurationMs()),
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
