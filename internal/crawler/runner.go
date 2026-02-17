package crawler

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
	"github.com/Dercraker/SearchEngine/internal/shared/customErrors"
	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
	"github.com/google/uuid"
)

type Runner struct {
	Logger     *slog.Logger
	SeedSource seeds.Source
	Processor  UrlProcessor

	CanonicalOptions seeds.CanonicalOptions
}

func (r Runner) RunOnce(ctx context.Context) (Stats, error) {
	rid := uuid.NewString()
	ctx = requestId.WithRunId(ctx, rid)

	st := Stats{StartTime: time.Now()}

	raw, err := r.SeedSource.Load(ctx)
	if err != nil {
		st.EndTime = time.Now()
		return st, err
	}

	list := seeds.SplitSeeds(raw)
	st.TotalSeeds = len(list)
	if len(list) == 0 {
		r.Logger.Error("[Crawler] no seeds provided")
		st.EndTime = time.Now()
		return st, err
	}

	seen := make(map[string]struct{}, len(list))

	for _, s := range list {
		u, nerr := seeds.NormalizeHTTPURL(s)
		if nerr != nil {
			st.InvalidSeeds++
			r.Logger.Error("[Crawler] invalid seed skipped", slog.Any("seed", s), slog.Any("error", nerr))
			continue
		}

		key, kerr := seeds.CanonicalKey(u, r.CanonicalOptions)
		if kerr != nil {
			st.InvalidSeeds++
			r.Logger.Error("[Crawler] invalid seed skipped", slog.Any("seed", s), slog.Any("error", kerr))
			continue
		}
		if _, ok := seen[key]; ok {
			st.DedupSkipped++
			continue
		}

		seen[key] = struct{}{}

		r.Logger.Info(string(obs.RunStart), slog.String("request_id", rid), slog.Int("seeds_count", len(list)), slog.String("seed", key), slog.String("url", u.String()), slog.String("canonical_key", key))
		cu, _ := url.Parse(key)

		st.Processed++

		if perr := r.Processor.Process(ctx, cu); perr != nil {
			if errors.Is(perr, customErrors.ErrMaxPagesReached) {
				r.Logger.Info("[Crawler] stop reason=max_pages_reached processed", slog.String("seed", key), slog.Int("processed", st.Processed))
				break
			}

			st.Failed++
			continue
		}
		st.Success++
	}
	st.EndTime = time.Now()
	r.Logger.Info("[Crawler] finished running")
	r.Logger.Info("[Crawler] Summary", slog.Int("totalSeed", st.TotalSeeds), slog.Int("InvalidSeeds", st.InvalidSeeds), slog.Int("DedupSkipped", st.DedupSkipped), slog.Int("processed", st.Processed), slog.Int("success", st.Success), slog.Int("failed", st.Failed), slog.Float64("DurationMs", st.DurationMs()))
	return st, nil
}
