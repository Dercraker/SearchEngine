package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Dercraker/SearchEngine/internal/DAL"
)

type QueueStore struct {
	Q *DAL.Queries
}

func (s QueueStore) Enqueue(ctx context.Context, url string) error {
	return s.Q.EnqueueURL(ctx, url)
}

func (s QueueStore) ClaimNextBatch(ctx context.Context, n int32) ([]DAL.ClaimNextBatchRow, error) {
	return s.Q.ClaimNextBatch(ctx, n)
}

func (s QueueStore) Ensure(ctx context.Context, url string) error {
	return s.Q.EnsureQueueURL(ctx, url)
}

func (s QueueStore) MarkCrawled(ctx context.Context, url string) error {
	return s.Q.MarkQueueCrawled(ctx, url)
}

func (s QueueStore) MarkFailed(ctx context.Context, url string, lastErr string, nextRunAt time.Time) error {
	return s.Q.MarkQueueFailed(ctx, DAL.MarkQueueFailedParams{
		Url: url,
		LastError: sql.NullString{
			String: lastErr,
			Valid:  true,
		},
		NextRunAt: nextRunAt,
	})
}

func (s QueueStore) ReleaseStale(ctx context.Context, staleAfter time.Duration) error {
	return s.Q.ReleaseStaleProcessing(ctx, sql.NullString{
		String: toPGInterval(staleAfter),
		Valid:  true,
	})
}

func toPGInterval(d time.Duration) string {
	sec := int(d.Round(time.Second).Seconds())
	if sec <= 0 {
		sec = 1
	}
	return fmt.Sprintf("%d seconds", sec)
}
