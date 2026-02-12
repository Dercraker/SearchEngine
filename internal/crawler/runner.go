package crawler

import (
	"context"
	"log/slog"
	"net/url"
	"time"

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
		if perr := r.Processor.Process(ctx, cu); perr != nil {
			st.Failed++
			r.Logger.Error("[Crawler] failed to process seed", slog.Any("seed", key), slog.Any("error", perr))
			continue
		}
		st.Success++
		r.Logger.Info("[Crawler] finished for seed", slog.Any("seed", key))
	}
	st.EndTime = time.Now()
	r.Logger.Info("[Crawler] finished running")
	r.Logger.Info("[Crawler] Summary", slog.Int("totalSeed", st.TotalSeeds), slog.Int("InvalidSeeds", st.InvalidSeeds), slog.Int("DedupSkipped", st.DedupSkipped), slog.Int("processed", st.Processed), slog.Int("success", st.Success), slog.Int("failed", st.Failed), slog.Duration("Duration", st.Duration()))
	return st, nil
}
