package config

import (
	"errors"

	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type Config struct {
	SeedFilePath string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	sfp := configHelper.GetEnv("API_ADDR", "")
	if sfp == "" {
		return Config{}, errors.New("SeedFilePath is required")
	}

	return Config{
		SeedFilePath: sfp,
	}, nil
}
