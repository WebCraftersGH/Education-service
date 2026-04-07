package middleware

import (
	"net/http"
	"time"

	"github.com/WebCraftersGH/Education-service/pgk/logging"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()
		startedAt := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		logger.WithFields(map[string]any{
			"method":      r.Method,
			"path":        r.URL.Path,
			"query":       r.URL.RawQuery,
			"status":      recorder.status,
			"remote_addr": r.RemoteAddr,
			"duration":    time.Since(startedAt).String(),
		}).Info("http request handled")
	})
}
