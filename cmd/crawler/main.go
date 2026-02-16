package main

import (
	"context"
	"log"
	"log/slog"
	"net/url"
	"os"

	"github.com/Dercraker/SearchEngine/internal/crawler"
	"github.com/Dercraker/SearchEngine/internal/crawler/config"
	"github.com/Dercraker/SearchEngine/internal/shared/logging"
)

type NoopProcessor struct{}

func (NoopProcessor) Process(_ context.Context, _ *url.URL) error {
	return nil
}

func main() {
	logger := logging.New()

	crawlerConfig, err := Config.LoadCrawlerConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	r := crawler.Buildcrawler(logger, crawlerConfig)

	ctx, cancel := context.WithTimeout(context.Background(), crawlerConfig.FetcherConfig.Timeout)
	defer cancel()

	stats, err := r.RunOnce(ctx)
	if err != nil {
		logger.Error("[Crawler]: failed to run once", slog.Any("error", err), slog.Any("stats", stats))
		os.Exit(1)
	}

}
