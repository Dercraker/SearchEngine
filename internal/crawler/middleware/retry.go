package middleware

import (
	"context"
	"net/url"
	"time"
)

func Retry(max int, backoff time.Duration) Middleware {
	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			var err error
			for i := 0; i <= max; i++ {
				err = next.Process(ctx, u)
				if err == nil {
					return nil
				}
				if i < max && backoff > 0 {
					select {
					case <-time.After(backoff):
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}
			return err
		})
	}
}
