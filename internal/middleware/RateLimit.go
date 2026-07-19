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

func RateLimit(RateLimit map[string]RequestTracker, TimeLimit map[time.Time][]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" || origin == "null" {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			tracker, ok := RateLimit[origin]
			if !ok {
				RateLimit[origin] = RequestTracker{amount: 1, firstRequest: time.Now(), lastRequest: time.Now()}
				TimeLimit[time.Now()] = append(TimeLimit[time.Now()], origin)
			} else {
				if time.Since(tracker.lastRequest) < time.Second && tracker.amount >= 4 {
					http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
					return
				}
				tracker.amount++
				slice := TimeLimit[tracker.lastRequest]

				newSlice := make([]string, 0, len(slice))
				for _, v := range slice {
					if v != origin {
						newSlice = append(newSlice, v)
					}
				}
				TimeLimit[tracker.lastRequest] = newSlice
				tracker.lastRequest = time.Now()
				RateLimit[origin] = tracker
				TimeLimit[tracker.lastRequest] = append(TimeLimit[time.Now()], origin)
			}

			for key, value := range TimeLimit {
				if time.Since(key) > 20*time.Minute {
					delete(TimeLimit, key)
					for _, v := range value {
						delete(RateLimit, v)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
