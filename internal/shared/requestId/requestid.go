package requestId

import "context"

type ctxKey struct{}

func WithRunId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

func GetRunId(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxKey{})
	s, isValid := v.(string)
	return s, isValid && s != ""
}
