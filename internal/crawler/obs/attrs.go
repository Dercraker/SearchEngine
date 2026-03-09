package obs

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/Dercraker/SearchEngine/internal/shared/instanceId"
	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
)

func BaseAttrs(ctx context.Context, u *url.URL) []any {
	attrs := make([]any, 0, 4)

	if rid, ok := requestId.GetRunId(ctx); ok {
		attrs = append(attrs, slog.String("request_id", rid))
	}
	if iid, ok := instanceId.GetInstanceId(ctx); ok {
		attrs = append(attrs, slog.String("instance_id", iid))
	}

	if u != nil {
		attrs = append(attrs,
			slog.String("url", u.String()),
			slog.String("host", u.Host),
		)
	}
	return attrs
}
