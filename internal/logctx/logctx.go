package logctx

import (
	"context"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"github.com/WebCraftersGH/Education-service/internal/requestctx"
)

const (
	requestIDKey = "request-id"
)

func WithContext(ctx context.Context, logger logging.Logger) logging.Logger {
	requestID, ok := requestctx.RequestID(ctx)
	if !ok || requestID == "" {
		return logger
	}

	return logger.WithField(requestIDKey, requestID)
}
