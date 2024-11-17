package middlewares

import (
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
)

// MetricsMiddleware is a middleware that collects and records various metrics
// about HTTP requests and responses. It tracks the number of active requests,
// the duration of each request, the size of the request body, and the size of
// the response body.
//
// It uses the following utility functions from the utils package:
// - IncrementActiveRequests: increments the count of active requests for a given method and path.
// - DecrementActiveRequests: decrements the count of active requests for a given method and path.
// - ObserveHTTPRequestDuration: records the duration of a completed request for a given method and path.
// - ObserveRequestSize: records the size of the request body for a given method and path.
// - ObserveResponseSize: records the size of the response body for a given method and path.
//
// The middleware wraps the response writer to capture the size of the response body.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		utils.IncrementActiveRequests(r.Method, r.URL.Path)

		defer func() {
			utils.DecrementActiveRequests(r.Method, r.URL.Path)
			duration := time.Since(start).Seconds()
			utils.ObserveHTTPRequestDuration(r.Method, r.URL.Path, duration)
		}()

		if r.ContentLength > 0 {
			utils.ObserveRequestSize(r.Method, r.URL.Path, float64(r.ContentLength))
		}

		wrappedWriter := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(wrappedWriter, r)
		utils.ObserveResponseSize(r.Method, r.URL.Path, float64(wrappedWriter.size))
	})
}

// responseWriter wraps http.ResponseWriter to capture response size
type responseWriter struct {
	http.ResponseWriter
	size int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
