// internal/crawler/middleware/outcome.go
package middleware

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/crawler/processors"
)

func OutcomeMW(logger *slog.Logger, qs processors.QueueStore, retryAfter string) Middleware {
	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			if qs != nil {
				_ = qs.Ensure(ctx, u.String())
			}

			err := next.Process(ctx, u)

			if qs == nil {
				return err
			}

			if err != nil {
				_ = qs.MarkFailed(ctx, u.String(), err.Error(), retryAfter)
				logger.Error(string(obs.QueueUpsert),
					append(obs.BaseAttrs(ctx, u),
						slog.String("queue_status", "failed"),
					)...,
				)
				return err
			}

			_ = qs.MarkCrawled(ctx, u.String())
			logger.Info(string(obs.QueueUpsert),
				append(obs.BaseAttrs(ctx, u),
					slog.String("queue_status", "crawled"),
				)...,
			)
			return nil
		})
	}
}
