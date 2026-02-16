package storage

import (
	"context"
	"database/sql"

	"github.com/Dercraker/SearchEngine/internal/DAL"
)

type DocumentStore struct {
	Q *DAL.Queries
}

func (s DocumentStore) UpsertFetch(ctx context.Context, url string, statusCode int, contentType string, contentHash string, raw []byte) error {
	_, err := s.Q.UpsertDocumentFetch(ctx, DAL.UpsertDocumentFetchParams{
		Url:         url,
		StatusCode:  int32(statusCode),
		ContentType: contentType,
		ContentHash: sql.NullString{
			String: contentHash,
			Valid:  true,
		},
		RawContent: raw,
	})
	return err
}
