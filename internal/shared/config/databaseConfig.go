package sharedconfig

import (
	"errors"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DatabaseDSN       string
	DBPingTimeout     time.Duration
	DBFailFast        bool
	DBMaxIdleConns    int
	DBMaxOpenConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration
}

func LoadDatabaseConfig() (DatabaseConfig, error) {
	_ = godotenv.Load()

	dsn := configHelper.GetEnv("DATABASE_DSN", "")
	if dsn == "" {
		return DatabaseConfig{}, errors.New("DATABASE_DSN is required")
	}

	dbPingTimeout, err := configHelper.ParseDuration("DATABASE_PING_TIMEOUT", "2s")
	if err != nil {
		return DatabaseConfig{}, err
	}

	dbFailFast := configHelper.GetEnvBool("DATABASE_FAIL_FAST", true)

	dbMaxIdleConns := configHelper.GetEnvInt("DATABASE_MAX_IDLE_CONNS", 10)
	dbMaxOpenConns := configHelper.GetEnvInt("DATABASE_MAX_OPEN_CONNS", 10)

	dbConnMaxLifetime, err := configHelper.ParseDuration("DATABASE_CONN_MAX_LIFETIME", "30m")
	if err != nil {
		return DatabaseConfig{}, err
	}

	dbConnMaxIdleTime, err := configHelper.ParseDuration("DATABASE_CONN_MAX_IDLE_TIME", "5m")
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		DatabaseDSN:       dsn,
		DBPingTimeout:     dbPingTimeout,
		DBFailFast:        dbFailFast,
		DBMaxIdleConns:    dbMaxIdleConns,
		DBMaxOpenConns:    dbMaxOpenConns,
		DBConnMaxLifetime: dbConnMaxLifetime,
		DBConnMaxIdleTime: dbConnMaxIdleTime,
	}, nil
}
