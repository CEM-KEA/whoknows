package middlewares

import (
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
)


// NoCacheMiddleware is a middleware that sets HTTP headers to prevent caching.
// It adds the following headers to the response:
// - Cache-Control: no-store, no-cache, must-revalidate, proxy-revalidate
// - Pragma: no-cache
// - Expires: 0
//
// This middleware logs the initialization, application, and completion of the no-cache headers.
//
// Parameters:
// - next: The next http.Handler in the middleware chain.
//
// Returns:
// - http.Handler: A handler that applies the no-cache headers and then calls the next handler.
func NoCacheMiddleware(next http.Handler) http.Handler {
	utils.LogInfo("Initializing no-cache middleware", nil)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.LogInfo("Applying no-cache headers", logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		})

		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		utils.LogInfo("No-cache headers applied", logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		})
		next.ServeHTTP(w, r)
	})
}
