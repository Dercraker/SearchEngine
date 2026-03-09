package instanceId

import "context"

type contextKey string

const instanceIdKey contextKey = "instance_Id"

func WithInstanceId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, instanceIdKey, id)
}

func GetInstanceId(ctx context.Context) (string, bool) {
	v := ctx.Value(instanceIdKey)
	s, isValid := v.(string)
	return s, isValid && s != ""
}
