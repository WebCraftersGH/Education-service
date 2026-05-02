package middleware

import (
	"github.com/google/uuid"
	"net/http"

	"github.com/WebCraftersGH/Education-service/internal/requestctx"
)

// GenerateRequestID generate request id
func GenerateRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		ctx := requestctx.WithRequestID(r.Context(), reqID)

		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}
