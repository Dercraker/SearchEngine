package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
	addr := getenv("API_ADDR", ":8080")

	rt, err := parseDuration("API_READ_TIMEOUT", "5s")
	if err != nil {
		return Config{}, err
	}

	wt, err := parseDuration("API_WRITE_TIMEOUT", "10s")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Addr:         addr,
		ReadTimeout:  rt,
		WriteTimeout: wt,
	}, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(key, fallback string) (time.Duration, error) {
	raw := getenv(key, fallback)
	d, err := time.ParseDuration(raw)

	if err != nil {
		return 0, errors.New("invalid duration: " + key + "=" + raw)
	}

	return d, nil
}
