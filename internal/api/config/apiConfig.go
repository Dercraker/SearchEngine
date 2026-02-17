package config

import (
	"errors"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/config"
	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	DatabaseConfig sharedconfig.DatabaseConfig

	SearchLimitDefault int
	SearchLimitMax     int
}

func LoadApiConfig() (ApiConfig, error) {
	_ = godotenv.Load()

	addr := configHelper.GetEnv("API_ADDR", ":8080")
	rt, err := configHelper.ParseDuration("API_READ_TIMEOUT", "5s")
	if err != nil {
		return ApiConfig{}, err
	}

	wt, err := configHelper.ParseDuration("API_WRITE_TIMEOUT", "10s")
	if err != nil {
		return ApiConfig{}, err
	}

	sld := configHelper.GetEnvInt("API_SEARCH_LIMIT_DEFAULT", 10)
	slm := configHelper.GetEnvInt("API_SEARCH_LIMIT_MAX", 50)

	dbConfig, err := sharedconfig.LoadDatabaseConfig()
	if err != nil {
		return ApiConfig{}, err
	}

	if sld <= 0 {
		sld = 10
	}
	if slm <= 0 {
		slm = 50
	}

	if sld > slm {
		return ApiConfig{}, errors.New("API_SEARCH_LIMIT_DEFAULT must be less than API_SEARCH_LIMIT_MAX")
	}

	return ApiConfig{
		Addr:         addr,
		ReadTimeout:  rt,
		WriteTimeout: wt,

		DatabaseConfig: dbConfig,

		SearchLimitDefault: sld,
		SearchLimitMax:     slm,
	}, nil
}
