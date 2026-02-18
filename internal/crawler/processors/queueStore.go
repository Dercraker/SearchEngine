package processors

import "context"

type QueueStore interface {
	Ensure(ctx context.Context, url string) error
	MarkCrawled(ctx context.Context, url string) error
	MarkFailed(ctx context.Context, url string, lastErr string, retryAfter string) error
}
