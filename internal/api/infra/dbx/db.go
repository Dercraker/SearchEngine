package dbx

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Options struct {
	DSN             string
	PingTimeout     time.Duration
	FailFast        bool
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Open(logger *slog.Logger, options Options) (*sql.DB, error) {
	if options.DSN == "" {
		return nil, errors.New("db: DSN is required")
	}

	db, err := sql.Open("pgx", options.DSN)
	if err != nil {
		logger.Error("db: failed to open connection", slog.Any("error", err))
		return nil, err
	}

	db.SetMaxOpenConns(options.MaxOpenConns)
	db.SetMaxIdleConns(options.MaxIdleConns)
	db.SetConnMaxIdleTime(options.ConnMaxIdleTime)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)

	if options.FailFast {
		ctx, cancel := context.WithTimeout(context.Background(), options.PingTimeout)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			logger.Error("db: failed to ping database", slog.Any("error", err))
			_ = db.Close()
			return nil, err
		}

		logger.Info("db: successfully connected to database")
	}

	return db, nil
}
