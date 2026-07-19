package middleware

import (
	"net/http"
	"time"
)

type RequestTracker struct {
	amount       int
	firstRequest time.Time
	lastRequest  time.Time
}

func RateLimit(RateLimit map[string]RequestTracker, TimeLimit map[time.Time]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}
