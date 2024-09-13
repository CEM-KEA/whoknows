package api

import (
	"github.com/CEM-KEA/whoknows/backend/internal/api/handlers"
	"github.com/CEM-KEA/whoknows/backend/internal/api/middlewares"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// Apply the CORS middleware
	router.Use(middlewares.CORSMiddleware)

	router.HandleFunc("/api/search", handlers.Search).Methods("POST")
	router.HandleFunc("/api/weather", nil).Methods("GET")   // Add the weather handler here
	router.HandleFunc("/api/register", nil).Methods("POST") // Add the register handler here
	router.HandleFunc("/api/login", nil).Methods("POST")    // Add the login handler here
	router.HandleFunc("/api/logout", nil).Methods("GET")    // Add the logout handler here

	return router
}
