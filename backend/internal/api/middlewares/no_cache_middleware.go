package middlewares

import (
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
)

// NoCacheMiddleware prevents caching by setting appropriate HTTP headers.
// Headers added:
// - Cache-Control: no-store, no-cache, must-revalidate, proxy-revalidate
// - Pragma: no-cache
// - Expires: 0
//
// Logs:
// - Logs initialization of the middleware during server startup.
// - Logs application of no-cache headers for each request.
//
// Parameters:
// - next: The next http.Handler in the middleware chain.
//
// Returns:
// - http.Handler: A handler that applies no-cache headers and then delegates to the next handler.
func NoCacheMiddleware(next http.Handler) http.Handler {
	utils.LogInfo("Initializing no-cache middleware", nil)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := utils.SanitizeValue(r.URL.Path)
		method := utils.SanitizeValue(r.Method)

		utils.LogInfo("Applying no-cache headers", map[string]interface{}{
			"path":   path,
			"method": method,
		})

		// Set headers to disable caching
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		utils.LogInfo("No-cache headers applied", map[string]interface{}{
			"path":   path,
			"method": method,
		})

		next.ServeHTTP(w, r)
	})
}
