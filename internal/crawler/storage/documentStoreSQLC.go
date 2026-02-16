package storage

import (
	"context"
	"database/sql"

	"github.com/Dercraker/SearchEngine/internal/DAL"
)

type DocumentStore struct {
	Q *DAL.Queries
}

func (s DocumentStore) GetHashByURL(ctx context.Context, url string) (string, error) {
	h, err := s.Q.GetDocumentHashByURL(ctx, url)
	if err != nil {
		return "", err
	}
	return h.String, err
}

func (s DocumentStore) TouchFetchAt(ctx context.Context, url string, statusCode int, contentType string) error {
	return s.Q.TouchDocumentFetchedAt(ctx, DAL.TouchDocumentFetchedAtParams{
		Url:         url,
		StatusCode:  int32(statusCode),
		ContentType: contentType,
	})
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
