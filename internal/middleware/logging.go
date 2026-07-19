package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := RequestIDFromContext(r.Context())
		if reqID == "" {
			reqID = "unknown"
		}
		rec := statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(&rec, r)

		log.Printf("[%s] %q;  STATUS: %v, RESPONSE TIME: %v (with id %s)\n ", r.Method, r.URL.Path, rec.status, time.Since(start), reqID)
	})
}
