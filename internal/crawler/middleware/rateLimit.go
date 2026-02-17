package middleware

import (
	"context"
	"errors"
	"net/url"

	"github.com/Dercraker/SearchEngine/internal/crawler/rateLimit"
)

func RateLimitMW(l *rateLimit.Limiter) Middleware {
	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			n := l.Count.Add(1)
			if n > l.Cfg.MaxPagesPerRun {
				return errors.ErrUnsupported
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
