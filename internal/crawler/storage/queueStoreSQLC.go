package storage

import (
	"context"
	"database/sql"
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
			Valid:  false,
		},
		NextRunAt: nextRunAt,
	})
}

func (s QueueStore) ReleaseStale(ctx context.Context, staleAfter time.Duration) error {
	return s.Q.ReleaseStaleProcessing(ctx, toPGInterval(staleAfter))
}

func toPGInterval(d time.Duration) sql.NullString {
	sec := int(d.Round(time.Second).Seconds())
	if sec <= 0 {
		return sql.NullString{
			String: "1",
			Valid:  true,
		}
	}
	return sql.NullString{
		String: (time.Duration(sec) * time.Second).String(),
		Valid:  true,
	}
}
