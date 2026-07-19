package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"time"
)

type contextKeyRequestid string

const requestIDKey contextKeyRequestid = "requestID"

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		seed := rand.NewSource(time.Now().UnixNano())
		n := rand.New(seed)

		length := 32
		result := make([]byte, length)
		for i := range result {
			result[i] = charset[n.Intn(len(charset))]
		}
		reqID := string(result)

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
