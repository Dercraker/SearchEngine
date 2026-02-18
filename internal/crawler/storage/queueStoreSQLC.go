package storage

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Dercraker/SearchEngine/internal/DAL"
)

type QueueStore struct {
	Q *DAL.Queries
}

func (s QueueStore) Ensure(ctx context.Context, url string) error {
	return s.Q.EnsureQueueURL(ctx, url)
}

func (s QueueStore) MarkCrawled(ctx context.Context, url string) error {
	return s.Q.MarkQueueCrawled(ctx, url)
}

func (s QueueStore) MarkFailed(ctx context.Context, url string, lastErr string, retryAfter string) error {
	retry, _ := strconv.ParseInt(retryAfter, 10, 64)
	return s.Q.MarkQueueFailed(ctx, DAL.MarkQueueFailedParams{
		Url: url,
		LastError: sql.NullString{
			String: lastErr,
			Valid:  false,
		},
		Column3: retry,
	})
}
