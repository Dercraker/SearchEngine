package crawler

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"time"

	errors2 "github.com/Dercraker/SearchEngine/internal/crawler/errors"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
)

type Runner struct {
	Logger     *slog.Logger
	SeedSource seeds.Source
	Processor  UrlProcessor

	CanonicalOptions seeds.CanonicalOptions
}

func (r Runner) RunOnce(ctx context.Context) (Stats, error) {
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
			r.Logger.Error("[Crawler] invalid seed skipped", slog.Any("seed", s), slog.Any("error", err))
		}
		if _, ok := seen[key]; ok {
			st.DedupSkipped++
			continue
		}

		seen[key] = struct{}{}

		r.Logger.Info("[Crawler] start for seed", slog.Any("seed", key))
		cu, _ := url.Parse(key)

		st.Processed++

		perr := r.Processor.Process(ctx, cu)
		if perr != nil {
			if errors.Is(err, errors2.ErrMaxPagesReached) {
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
