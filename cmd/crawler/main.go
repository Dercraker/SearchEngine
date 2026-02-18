package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Dercraker/SearchEngine/internal/crawler"
	"github.com/Dercraker/SearchEngine/internal/crawler/config"
	"github.com/Dercraker/SearchEngine/internal/shared/logging"
)

func main() {
	logger := logging.New()

	crawlerConfig, err := config.LoadCrawlerConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	runner := crawler.BuildCrawler(logger, crawlerConfig)

	ctx, cancel := context.WithTimeout(context.Background(), crawlerConfig.RunTimeout)
	defer cancel()

	stats, err := runner.RunOnce(ctx)
	if err != nil {
		logger.Error("[Crawler]: failed to run once", slog.Any("error", err), slog.Any("stats", stats))
		os.Exit(1)
	}

}
