package main

import (
	"context"
	"log"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler"
	"github.com/Dercraker/SearchEngine/internal/crawler/config"
	"github.com/Dercraker/SearchEngine/internal/crawler/middleware"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
	"github.com/Dercraker/SearchEngine/internal/shared/logging"
)

type NoopProcessor struct{}

func (NoopProcessor) Process(_ context.Context, _ *url.URL) error {
	return nil
}

func main() {
	logger := logging.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	p := NoopProcessor{}

	proc := middleware.Chain(
		p,
		middleware.LoggingMW(logger),
		middleware.Retry(2, 250*time.Millisecond),
	)

	r := crawler.Runner{
		Logger:     logger,
		Processor:  proc,
		SeedSource: seeds.FileSource{Path: cfg.SeedFilePath},
		CanonicalOptions: seeds.CanonicalOptions{
			DropTrackingParams: true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats, err := r.RunOnce(ctx)
	if err != nil {
		logger.Error("[Crawler]: failed to run once", slog.Any("error", err), slog.Any("stats", stats))
		os.Exit(1)
	}

}
