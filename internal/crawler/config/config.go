package CrawlerConfig

import (
	"errors"

	"github.com/Dercraker/SearchEngine/internal/crawler/httpfetch"
	"github.com/Dercraker/SearchEngine/internal/shared/config"
	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type Config struct {
	SeedFilePath   string
	FetcherConfig  httpfetch.Config
	DatabaseConfig config.DatabaseConfig
}

func Load() (Config, error) {
	_ = godotenv.Load()

	sfp := configHelper.GetEnv("CRAWLER_SEED_FILE_PATH", "")
	if sfp == "" {
		return Config{}, errors.New("SeedFilePath is required")
	}

	fetchConfig, err := httpfetch.LoadFetcherConfig()
	if err != nil {
		return Config{}, err
	}

	dbConfig, err := config.LoadDatabaseConfig()

	return Config{
		SeedFilePath:   sfp,
		FetcherConfig:  fetchConfig,
		DatabaseConfig: dbConfig,
	}, nil
}
