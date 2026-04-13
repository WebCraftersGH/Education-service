package requestctx

import (
	"context"
)

const (
	requestIDKey = "request-id"
)

func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}
