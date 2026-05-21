package middleware

import (
	"net/http"
	"time"

	"github.com/WebCraftersGH/Education-service/internal/requestctx"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}

	n, err := r.ResponseWriter.Write(data)
	r.bytes += n
	return n, err
}

func RequestLogger(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			recorder := &statusRecorder{ResponseWriter: w}

			requestID, _ := requestctx.RequestID(r.Context())
			logger.WithFields(map[string]any{
				"request_id":     requestID,
				"method":         r.Method,
				"path":           r.URL.Path,
				"raw_query":      r.URL.RawQuery,
				"remote_addr":    r.RemoteAddr,
				"user_agent":     r.UserAgent(),
				"content_length": r.ContentLength,
			}).Debug("http request started")

			next.ServeHTTP(recorder, r)

			if recorder.status == 0 {
				recorder.status = http.StatusOK
			}

			logger.WithFields(map[string]any{
				"request_id":     requestID,
				"method":         r.Method,
				"path":           r.URL.Path,
				"raw_query":      r.URL.RawQuery,
				"status":         recorder.status,
				"response_bytes": recorder.bytes,
				"duration_ms":    time.Since(startedAt).Milliseconds(),
			}).Debug("http request completed")
		})
	}
}
