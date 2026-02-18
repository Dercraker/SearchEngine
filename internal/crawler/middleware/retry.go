package middleware

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/url"
	"syscall"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
	"github.com/Dercraker/SearchEngine/internal/shared/customErrors"
)

func Retry(logger *slog.Logger, st *obs.Stats, max int, backoff time.Duration) Middleware {
	return func(next URLProcessor) URLProcessor {
		return URLProcessorFunc(func(ctx context.Context, u *url.URL) error {
			var err error

			for attempt := 0; attempt <= max; attempt++ {
				err = next.Process(ctx, u)
				if err == nil {
					return nil
				}

				if attempt == max || !shouldRetry(err) {
					return err
				}
				if st != nil {
					st.Retries.Add(1)
				}

				logger.Warn(string(obs.URLRetry),
					append(obs.BaseAttrs(ctx, u),
						slog.Int("attempt", attempt),
						slog.Int("max", max),
						slog.Any("error", err),
						slog.Duration("backoff", backoff),
					)...,
				)

				select {
				case <-time.After(backoff):
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return err
		})
	}
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, customErrors.ErrBodyTooLarge) {
		return false
	}
	if errors.Is(err, customErrors.ErrTooManyRedirects) {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return false
	}

	// Retry sur timeout
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// net.Error (timeout/temporary)
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout() || ne.Temporary()
	}

	// EOF / reset / broken pipe (classiques)
	if errors.Is(err, io.EOF) {
		return true
	}
	if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.EPIPE) {
		return true
	}

	return false
}
