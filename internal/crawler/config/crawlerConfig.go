package config

import (
	"errors"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/config"
	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type CrawlerConfig struct {
	SeedFilePath  string
	RunTimeout    time.Duration
	FetcherConfig FetcherConfig
	LimitConfig   LimitConfig

	DatabaseConfig sharedconfig.DatabaseConfig
}

func LoadCrawlerConfig() (CrawlerConfig, error) {
	_ = godotenv.Load()

	sfp := configHelper.GetEnv("CRAWLER_SEED_FILE_PATH", "")
	if sfp == "" {
		return CrawlerConfig{}, errors.New("SeedFilePath is required")
	}

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

	return CrawlerConfig{
		SeedFilePath:   sfp,
		RunTimeout:     runTimeout,
		FetcherConfig:  fetchConfig,
		LimitConfig:    limitConfig,
		DatabaseConfig: dbConfig,
	}, nil
}
