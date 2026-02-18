package crawler

import (
	"log/slog"
	"time"

	"github.com/Dercraker/SearchEngine/internal/DAL"
	"github.com/Dercraker/SearchEngine/internal/api/infra/dbx"
	CrawlerConfig "github.com/Dercraker/SearchEngine/internal/crawler/config"
	"github.com/Dercraker/SearchEngine/internal/crawler/httpfetch"
	"github.com/Dercraker/SearchEngine/internal/crawler/middleware"
	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/crawler/processors"
	"github.com/Dercraker/SearchEngine/internal/crawler/rateLimit"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
	"github.com/Dercraker/SearchEngine/internal/crawler/storage"
)

func BuildCrawler(logger *slog.Logger, cfg CrawlerConfig.CrawlerConfig) *QueueRunner {
	stats := &obs.Stats{}

	seedSource := seeds.FileSource{Path: cfg.SeedFilePath}

	dbConn, err := dbx.Open(logger, dbx.Options{
		DSN:             cfg.DatabaseConfig.DatabaseDSN,
		PingTimeout:     cfg.DatabaseConfig.DBPingTimeout,
		FailFast:        cfg.DatabaseConfig.DBFailFast,
		MaxIdleConns:    cfg.DatabaseConfig.DBMaxIdleConns,
		MaxOpenConns:    cfg.DatabaseConfig.DBMaxOpenConns,
		ConnMaxLifetime: cfg.DatabaseConfig.DBConnMaxLifetime,
		ConnMaxIdleTime: cfg.DatabaseConfig.DBConnMaxIdleTime,
	})

	if err != nil {
		panic(err)
	}

	q := DAL.New(dbConn)

	documentStore := storage.DocumentStore{Q: q}
	queueStore := storage.QueueStore{Q: q}

	fetcherCfg := cfg.FetcherConfig
	fetcherCfg.Logger = logger
	fetcher := httpfetch.New(fetcherCfg)

	rateLimiter := rateLimit.NewRateLimiter(cfg.LimitConfig)

	downloader := processors.Downloader{
		Fetcher: fetcher,
		Store:   documentStore,
		Stats:   stats,
	}

	proc := middleware.Chain(
		downloader,
		middleware.RateLimitMW(rateLimiter),
		middleware.Retry(logger, stats, 2, 250*time.Millisecond),
		middleware.OutcomeMW(logger, queueStore, "10s"),
		middleware.LoggingMW(logger),
	)

	return &QueueRunner{
		Logger:           logger,
		SeedSource:       seedSource,
		Processor:        proc,
		Queue:            queueStore,
		Stats:            stats,
		CanonicalOptions: seeds.CanonicalOptions{DropTrackingParams: true},
		batchSize:        50,
		StaleAfter:       10 * time.Minute,
		MaxPagesPerRun:   cfg.LimitConfig.MaxPagesPerRun,
	}
}
