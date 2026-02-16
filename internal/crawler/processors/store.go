package processors

import "context"

type DocumentStore interface {
	UpsertFetch(ctx context.Context, url string, statusCode int, contentType string, hash string, rawContent []byte) error
}
