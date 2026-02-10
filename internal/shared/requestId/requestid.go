package requestId

import "context"

type ctxKey struct{}

func With(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

func Get(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxKey{})
	s, isValid := v.(string)
	return s, isValid && s != ""
}
