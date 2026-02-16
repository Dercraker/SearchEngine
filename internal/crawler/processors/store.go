package processors

import "context"

type DocumentStore interface {
	GetHashByURL(ctx context.Context, url string) (string, error)
	TouchFetchAt(ctx context.Context, url string, statusCode int, contentType string) error
	UpsertFetch(ctx context.Context, url string, statusCode int, contentType string, hash string, rawContent []byte) error
}
