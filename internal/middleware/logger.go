package middleware

import (
	"net/http"
	"github.com/google/uuid"
	"context"

)

const requestIDKey = "request-id"

//GenerateRequestID generate request id
func GenerateRequestID (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)	
		
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)	
	})
}
