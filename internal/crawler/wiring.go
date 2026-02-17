package crawler

import (
	"log/slog"
	"time"

	"github.com/Dercraker/SearchEngine/internal/DAL"
	"github.com/Dercraker/SearchEngine/internal/api/infra/dbx"
	CrawlerConfig "github.com/Dercraker/SearchEngine/internal/crawler/config"
	"github.com/Dercraker/SearchEngine/internal/crawler/httpfetch"
	"github.com/Dercraker/SearchEngine/internal/crawler/middleware"
	"github.com/Dercraker/SearchEngine/internal/crawler/processors"
	"github.com/Dercraker/SearchEngine/internal/crawler/rateLimit"
	"github.com/Dercraker/SearchEngine/internal/crawler/seeds"
	"github.com/Dercraker/SearchEngine/internal/crawler/storage"
)

func BuildCrawler(logger *slog.Logger, cfg CrawlerConfig.CrawlerConfig) Runner {
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

	store := storage.DocumentStore{Q: q}

	fetcherCfg := cfg.FetcherConfig
	fetcherCfg.Logger = logger
	fetcher := httpfetch.New(fetcherCfg)
	downloader := processors.Downloader{Fetcher: fetcher, Store: store}

	limiter := rateLimit.NewRateLimiter(cfg.LimitConfig)

	proc := middleware.Chain(
		downloader,
		middleware.RateLimitMW(limiter),
		middleware.LoggingMW(logger),
		middleware.Retry(2, 250*time.Millisecond),
	)

	return Runner{
		SeedSource: seedSource,
		Logger:     logger,
		Processor:  proc,
		CanonicalOptions: seeds.CanonicalOptions{
			DropTrackingParams: true,
		},
	}
}
