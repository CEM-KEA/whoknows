package api

import (
	"net/http"

	_ "github.com/CEM-KEA/whoknows/backend/docs" // docs is generated by Swag CLI
	"github.com/CEM-KEA/whoknows/backend/internal/api/handlers"
	"github.com/CEM-KEA/whoknows/backend/internal/api/middlewares"
	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter sets up the application router
func NewRouter() http.Handler {
	utils.LogInfo("Initializing router", nil)

	router := mux.NewRouter()

	// Static file handlers
	utils.LogInfo("Setting up static file routes", nil)
	router.HandleFunc("/api/robots.txt", serveStaticFile("./static/robots.txt", "text/plain"))
	router.HandleFunc("/api/sitemap.xml", serveStaticFile("./static/sitemap.xml", "application/xml"))

	// Redirects
	utils.LogInfo("Setting up redirects", nil)
	setupRedirects(router)

	// Swagger documentation
	utils.LogInfo("Setting up Swagger documentation route", nil)
	router.PathPrefix("/api/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes
	utils.LogInfo("Setting up API routes", nil)
	setupAPIRoutes(router)

	// CORS configuration
	corsHandler := setupCORS()

	// Wrap router with middleware
	utils.LogInfo("Applying middlewares", nil)
	return middlewares.MetricsMiddleware(middlewares.NoCacheMiddleware(corsHandler(router)))
}

// serveStaticFile serves a static file with the specified content type
func serveStaticFile(filePath string, contentType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.LogInfo("Serving static file", logrus.Fields{
			"path": filePath,
		})
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, filePath)
	}
}

// setupRedirects sets up the necessary URL redirects
func setupRedirects(router *mux.Router) {
	redirectHandler := RedirectToSwaggerHandler()
	router.Handle("/", redirectHandler)
	router.Handle("/api", redirectHandler)
	router.Handle("/api/", redirectHandler)
	router.Handle("/api/swagger", redirectHandler)
	router.Handle("/robots.txt", http.RedirectHandler("/api/robots.txt", http.StatusMovedPermanently))
	router.Handle("/sitemap.xml", http.RedirectHandler("/api/sitemap.xml", http.StatusMovedPermanently))
	utils.LogInfo("Redirects configured", nil)
}

// setupAPIRoutes configures the API endpoints
func setupAPIRoutes(router *mux.Router) {
	router.HandleFunc("/api/search", handlers.Search).Methods("GET")
	router.HandleFunc("/api/weather", handlers.WeatherHandler).Methods("GET")
	router.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/login", handlers.Login).Methods("POST")
	router.HandleFunc("/api/logout", handlers.LogoutHandler).Methods("GET")
	router.HandleFunc("/api/validate-login", handlers.ValidateLoginHandler).Methods("GET")
	router.HandleFunc("/api/change-password", handlers.ChangePasswordHandler).Methods("POST")
	utils.LogInfo("API routes configured", nil)
}

// setupCORS sets up CORS configuration based on the environment
func setupCORS() func(http.Handler) http.Handler {
	env := config.AppConfig.Environment.Environment
	var allowedOrigins []string

	switch env {
	case "development":
		allowedOrigins = []string{"*"}
	case "test":
		allowedOrigins = []string{"http://localhost", "https://localhost"}
	default:
		allowedOrigins = []string{"http://cemdev.dk", "https://cemdev.dk"}
	}

	utils.LogInfo("Configuring CORS", logrus.Fields{
		"environment": env,
		"allowedOrigins": allowedOrigins,
	})

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler
}

// RedirectToSwaggerHandler handles redirection to Swagger documentation
func RedirectToSwaggerHandler() http.Handler {
	utils.LogInfo("Redirecting to Swagger documentation", nil)
	return http.RedirectHandler("/api/swagger/", http.StatusMovedPermanently)
}
