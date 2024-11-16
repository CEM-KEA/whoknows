package middlewares

import (
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
)

// NoCacheMiddleware prevents caching by setting appropriate headers
func NoCacheMiddleware(next http.Handler) http.Handler {
	utils.LogInfo("Initializing no-cache middleware", nil)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set no-cache headers
		utils.LogInfo("Applying no-cache headers", logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		})

		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Log and proceed to the next handler
		utils.LogInfo("No-cache headers applied", logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		})
		next.ServeHTTP(w, r)
	})
}
