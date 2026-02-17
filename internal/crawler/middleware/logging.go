package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/shared/customErrors"
)

func LoggingMW(logger *slog.Logger) Middleware {
	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			start := time.Now()

			logger.Info(string(obs.URLStart), obs.BaseAttrs(ctx, u)...)

			err := next.Process(ctx, u)
			d := time.Since(start)

			attrs := append(obs.BaseAttrs(ctx, u), slog.Duration("duration", d))

			if errors.Is(err, customErrors.ErrMaxPagesReached) {
				logger.Info(string(obs.URLEnd), append(attrs,
					slog.String("result", "stop_max_pages"),
				)...)
				return err
			}

			if err != nil {
				logger.Error(string(obs.URLEndFailed), append(attrs,
					slog.String("result", "failed"),
					slog.Any("error", err),
				)...)
				return err
			}

			logger.Info(string(obs.URLEnd), append(attrs,
				slog.String("result", "ok"),
			)...)

			return nil
		})
	}
}

type URLProcessorFunc func(ctx context.Context, u *url.URL) error

func (f URLProcessorFunc) Process(ctx context.Context, u *url.URL) error {
	return f(ctx, u)
}
