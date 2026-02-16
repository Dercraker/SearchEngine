package middleware

import (
	"context"
	"net/url"
	"sync"
	"time"

	Config "github.com/Dercraker/SearchEngine/internal/crawler/config"
	"github.com/Dercraker/SearchEngine/internal/crawler/errors"
)

func Throttle(limits Config.LimitConfig) Middleware {
	var mu sync.Mutex
	count := 0

	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			mu.Lock()
			defer mu.Unlock()

			if limits.MaxPagesPerRun > 0 && count >= limits.MaxPagesPerRun {
				mu.Unlock()
				return errors.ErrMaxPagesReached
			}
			mu.Unlock()

			err := next.Process(ctx, u)

			mu.Lock()
			count++
			mu.Unlock()

			if limits.DelayBetweenRequests > 0 {
				select {
				case <-time.After(limits.DelayBetweenRequests):
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			return err
		})
	}
}
