// File: internal/middleware/logging.go
package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		rw := &responseWriter{w, http.StatusOK}

		// Process request
		next.ServeHTTP(rw, r)

		// Log request details
		log.Printf(
			"Method: %s Path: %s Status: %d Duration: %v",
			r.Method,
			r.URL.Path,
			rw.status,
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// RequestLogger creates a middleware for request logging
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{w, http.StatusOK}

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Create a new context with the request ID
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		r = r.WithContext(ctx)

		// Add the request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		defer func() {
			log.Printf(
				"[%s] %s %s %d %v",
				requestID,
				r.Method,
				r.URL.Path,
				ww.status,
				time.Since(start),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}
