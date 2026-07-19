package middleware

import (
	"SocialMedia/internal/auth"
	"context"
	"net/http"
	"time"
)

type contextKey string

const userIDKey contextKey = "userID"

type SessionsStore interface {
	Get(ctx context.Context, sessionID string) (auth.Session, error)
}

func NewAuthMiddleware(sessions SessionsStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId, err := r.Cookie("sessionId")
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			session, err := sessions.Get(r.Context(), sessionId.Value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if session.ExpiresAt.Before(time.Now()) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			}

			ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
