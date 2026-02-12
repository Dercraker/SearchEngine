package middleware

import (
	"context"
	"log/slog"
	"net/url"
	"time"
)

func LoggingMW(logger *slog.Logger) Middleware {
	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			start := time.Now()
			logger.Info("[Crawler] Processing URL", slog.String("url", u.String()))
			err := next.Process(ctx, u)
			d := time.Since(start)

			if err != nil {
				logger.Error("[Crawler] Failed to process URL", slog.String("url", u.String()), slog.Duration("duration", d), slog.Any("error", err))
				return err
			}

			logger.Info("[Crawler] Successfully processed URL", slog.String("url", u.String()), slog.Duration("duration", d))
			return nil
		})
	}
}

type URLProcessorFunc func(ctx context.Context, u *url.URL) error

func (f URLProcessorFunc) Process(ctx context.Context, u *url.URL) error {
	return f(ctx, u)
}
