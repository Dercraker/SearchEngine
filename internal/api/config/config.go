package config

import (
	"errors"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
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

	dsn := configHelper.GetEnv("DATABASE_DSN", "")
	if dsn == "" {
		return Config{}, errors.New("DATABASE_DSN is required")
	}

	dbPingTimeout, err := configHelper.ParseDuration("DATABASE_PING_TIMEOUT", "2s")
	if err != nil {
		return Config{}, err
	}

	dbFailFast := configHelper.GetEnvBool("DATABASE_FAIL_FAST", true)

	dbMaxIdleConns := configHelper.GetEnvInt("DATABASE_MAX_IDLE_CONNS", 10)
	dbMaxOpenConns := configHelper.GetEnvInt("DATABASE_MAX_OPEN_CONNS", 10)

	dbConnMaxLifetime, err := configHelper.ParseDuration("DATABASE_CONN_MAX_LIFETIME", "30m")
	if err != nil {
		return Config{}, err
	}

	dbConnMaxIdleTime, err := configHelper.ParseDuration("DATABASE_CONN_MAX_IDLE_TIME", "5m")
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
