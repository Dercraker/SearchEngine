package middleware

import (
	"context"
	"net/url"
)

type URLProcessor interface {
	Process(ctx context.Context, u *url.URL) error
}

type Middleware func(next URLProcessor) URLProcessor

func Chain(p URLProcessor, middlewares ...Middleware) URLProcessor {
	for i := len(middlewares) - 1; i >= 0; i-- {
		p = middlewares[i](p)
	}
	return p
}
