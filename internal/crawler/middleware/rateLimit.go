package middleware

import (
	"context"
	"net/url"

	"github.com/Dercraker/SearchEngine/internal/crawler/rateLimit"
	"github.com/Dercraker/SearchEngine/internal/shared/customErrors"
)

func RateLimitMW(l *rateLimit.Limiter) Middleware {
	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			n := l.Count.Add(1)
			if n > l.Cfg.MaxPagesPerRun {
				return customErrors.ErrMaxPagesReached
			}

			if err := l.WaitGlobal(ctx); err != nil {
				return err
			}

			host := rateLimit.NormalizeHost(u.Host)
			hl := l.GetHostLimiter(host)
			if err := l.WaitHost(ctx, hl); err != nil {
				return err
			}

			return next.Process(ctx, u)
		})
	}
}
