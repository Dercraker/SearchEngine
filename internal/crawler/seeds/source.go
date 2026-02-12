package seeds

import "context"

type Source interface {
	Load(ctx context.Context) ([]string, error)
}
