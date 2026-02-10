package config

import (
	"errors"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	DatabaseDSN       string
	DBPingTimeout     time.Duration
	DBFailFast        bool
	DBMaxIdleConns    int
	DBMaxOpenConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	SearchLimitDefault int
	SearchLimitMax     int
}

func Load() (Config, error) {
	_ = godotenv.Load()

	addr := getEnv("API_ADDR", ":8080")
	rt, err := parseDuration("API_READ_TIMEOUT", "5s")
	if err != nil {
		return Config{}, err
	}

	wt, err := parseDuration("API_WRITE_TIMEOUT", "10s")
	if err != nil {
		return Config{}, err
	}

	sld := getEnvInt("API_SEARCH_LIMIT_DEFAULT", 10)
	slm := getEnvInt("API_SEARCH_LIMIT_MAX", 50)

	dsn := getEnv("DATABASE_DSN", "")
	if dsn == "" {
		return Config{}, errors.New("DATABASE_DSN is required")
	}

	dbPingTimeout, err := parseDuration("DATABASE_PING_TIMEOUT", "2s")
	if err != nil {
		return Config{}, err
	}

	dbFailFast := getEnvBool("DATABASE_FAIL_FAST", true)

	dbMaxIdleConns := getEnvInt("DATABASE_MAX_IDLE_CONNS", 10)
	dbMaxOpenConns := getEnvInt("DATABASE_MAX_OPEN_CONNS", 10)

	dbConnMaxLifetime, err := parseDuration("DATABASE_CONN_MAX_LIFETIME", "30m")
	if err != nil {
		return Config{}, err
	}

	dbConnMaxIdleTime, err := parseDuration("DATABASE_CONN_MAX_IDLE_TIME", "5m")
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

		DatabaseDSN:       dsn,
		DBPingTimeout:     dbPingTimeout,
		DBFailFast:        dbFailFast,
		DBMaxIdleConns:    dbMaxIdleConns,
		DBMaxOpenConns:    dbMaxOpenConns,
		DBConnMaxLifetime: dbConnMaxLifetime,
		DBConnMaxIdleTime: dbConnMaxIdleTime,

		SearchLimitDefault: sld,
		SearchLimitMax:     slm,
	}, nil
}
