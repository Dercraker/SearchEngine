package obs

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
)

func BaseAttrs(ctx context.Context, u *url.URL) []any {
	attrs := []any{
		slog.String("url", u.String()),
		slog.String("host", u.Host),
	}

	if rid, ok := requestId.GetRunId(ctx); ok {
		attrs = append(attrs, slog.String("request_id", rid))
	}
	return attrs
}
