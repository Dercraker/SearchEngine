package Config

import (
	"errors"

	"github.com/Dercraker/SearchEngine/internal/shared/config"
	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type CrawlerConfig struct {
	SeedFilePath  string
	FetcherConfig FetcherConfig
	LimitConfig   LimitConfig

	DatabaseConfig SharedConfig.DatabaseConfig
}

func LoadCrawlerConfig() (CrawlerConfig, error) {
	_ = godotenv.Load()

	sfp := configHelper.GetEnv("CRAWLER_SEED_FILE_PATH", "")
	if sfp == "" {
		return CrawlerConfig{}, errors.New("SeedFilePath is required")
	}

	fetchConfig, err := LoadFetcherConfig()
	if err != nil {
		return CrawlerConfig{}, err
	}

	limitConfig, err := LoadLimitConfig()
	if err != nil {
		return CrawlerConfig{}, err
	}

	dbConfig, err := SharedConfig.LoadDatabaseConfig()

	return CrawlerConfig{
		SeedFilePath:   sfp,
		FetcherConfig:  fetchConfig,
		LimitConfig:    limitConfig,
		DatabaseConfig: dbConfig,
	}, nil
}
