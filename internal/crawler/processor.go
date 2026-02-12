package crawler

import (
	"context"
	"net/url"
)

type UrlProcessor interface {
	Process(ctx context.Context, url *url.URL) error
}
