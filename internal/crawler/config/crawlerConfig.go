package config

import (
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/config"
	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type CrawlerConfig struct {
	RunTimeout    time.Duration
	FetcherConfig FetcherConfig
	LimitConfig   LimitConfig
	WorkerConfig  WorkerConfig

	DatabaseConfig sharedconfig.DatabaseConfig
}

func LoadCrawlerConfig() (CrawlerConfig, error) {
	_ = godotenv.Load()

	runTimeout, err := configHelper.ParseDuration("CRAWLER_RUN_TIMEOUT", "1h")
	if err != nil {
		return CrawlerConfig{}, err
	}

	fetchConfig, err := LoadFetcherConfig()
	if err != nil {
		return CrawlerConfig{}, err
	}

	limitConfig, err := LoadLimitConfig()
	if err != nil {
		return CrawlerConfig{}, err
	}

	dbConfig, err := sharedconfig.LoadDatabaseConfig()
	if err != nil {
		return CrawlerConfig{}, err
	}

	workerConfig, err := LoadWorkerConfig()
	if err != nil {
		return CrawlerConfig{}, err
	}

	return CrawlerConfig{
		RunTimeout:    runTimeout,
		FetcherConfig: fetchConfig,
		LimitConfig:   limitConfig,
		WorkerConfig:  workerConfig,

		DatabaseConfig: dbConfig,
	}, nil
}
