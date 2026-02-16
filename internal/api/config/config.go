package config

import (
	"errors"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/config"
	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	DatabaseConfig config.DatabaseConfig

	SearchLimitDefault int
	SearchLimitMax     int
}

func Load() (Config, error) {
	_ = godotenv.Load()

	addr := configHelper.GetEnv("API_ADDR", ":8080")
	rt, err := configHelper.ParseDuration("API_READ_TIMEOUT", "5s")
	if err != nil {
		return Config{}, err
	}

	wt, err := configHelper.ParseDuration("API_WRITE_TIMEOUT", "10s")
	if err != nil {
		return Config{}, err
	}

	sld := configHelper.GetEnvInt("API_SEARCH_LIMIT_DEFAULT", 10)
	slm := configHelper.GetEnvInt("API_SEARCH_LIMIT_MAX", 50)

	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		return Config{}, err
	}

	if sld <= 0 {
		sld = 10
	}
	if slm <= 0 {
		slm = 50
	}

	if sld > slm {
		return Config{}, errors.New("API_SEARCH_LIMIT_DEFAULT must be less than API_SEARCH_LIMIT_MAX")
	}

	return Config{
		Addr:         addr,
		ReadTimeout:  rt,
		WriteTimeout: wt,

		DatabaseConfig: dbConfig,

		SearchLimitDefault: sld,
		SearchLimitMax:     slm,
	}, nil
}
