package middleware

import (
	"errors"
	"strings"
	"context"
	"net/http"	
	"github.com/google/uuid"

	"github.com/WebCraftersGH/Education-service/internal/authclient"
	"github.com/WebCraftersGH/Education-service/internal/requestctx"
)

type AuthChecker interface{
	Check(ctx context.Context, token string) (uuid.UUID, error)
}

func AuthFromCookie(cookieName string, authChecker AuthChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(cookie.Value)
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := authChecker.Check(r.Context(), token)
			if err != nil {
				if errors.Is(err, authclient.ErrUnauthorized) {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}

				http.Error(w, "auth service unavailable", http.StatusServiceUnavailable)
				return
			}

			ctx := requestctx.WithUserID(r.Context(), userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
