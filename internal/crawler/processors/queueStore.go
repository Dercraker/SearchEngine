package processors

import (
	"context"
	"time"

	"github.com/Dercraker/SearchEngine/internal/DAL"
)

type QueueStore interface {
	Ensure(ctx context.Context, url string) error
	Enqueue(ctx context.Context, url string) error
	ClaimNextBatch(ctx context.Context, n int32) ([]DAL.ClaimNextBatchRow, error)
	MarkCrawled(ctx context.Context, url string) error
	MarkFailed(ctx context.Context, url string, lastErr string, retryAfter time.Time) error
	ReleaseStale(ctx context.Context, staleAfter time.Duration) error
}
