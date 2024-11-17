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
		path := utils.SanitizeValue(r.URL.Path) // Ensure the path is sanitized for logging and metrics.
		method := utils.SanitizeValue(r.Method)

		utils.LogInfo("Incoming request", map[string]interface{}{
			"method": method,
			"path":   path,
		})

		utils.IncrementActiveRequests(method, path)
		
		wrappedWriter := &responseWriter{ResponseWriter: w}

		defer func() {
			utils.DecrementActiveRequests(method, path)

			duration := time.Since(start).Seconds()
			utils.ObserveHTTPRequestDuration(method, path, duration)

			utils.LogInfo("Request completed", map[string]interface{}{
				"method":       method,
				"path":         path,
				"duration_sec": duration,
				"response_size": wrappedWriter.size,
			})
		}()

		if r.ContentLength > 0 {
			utils.ObserveRequestSize(method, path, float64(r.ContentLength))
		}

		next.ServeHTTP(wrappedWriter, r)
		utils.ObserveResponseSize(method, path, float64(wrappedWriter.size))
	})
}

// responseWriter wraps http.ResponseWriter to capture response size.
type responseWriter struct {
	http.ResponseWriter
	size int
}

// Write intercepts the response body write to track the size.
func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
