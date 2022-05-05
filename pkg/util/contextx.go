package util

import "context"

// NewStrContext returns a new Context that carries value strVal.
func NewStrContext(ctx context.Context, key, strVal string) context.Context {
	return context.WithValue(ctx, key, strVal)
}

// FromContextForStr returns the string	 value stored in ctx, if any.
func FromContextForStr(ctx context.Context, key string) (string, bool) {
	u, ok := ctx.Value(key).(string)
	return u, ok
}
