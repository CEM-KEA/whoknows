package api

import (
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/api/handlers"
	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/search", handlers.Search).Methods("POST")
	router.HandleFunc("/api/weather", nil).Methods("GET")   // Add the weather handler here
	router.HandleFunc("/api/register", nil).Methods("POST") // Add the register handler here
	router.HandleFunc("/api/login", nil).Methods("POST")    // Add the login handler here
	router.HandleFunc("/api/logout", nil).Methods("GET")    // Add the logout handler here

	// if environment is not production, allow all origins (*)
	var allowedOrigins []string

	if config.AppConfig.Environment.Environment != "production" {
		allowedOrigins = []string{"*"}
	} else {
		allowedOrigins = []string{"http://frontend", "http://localhost:80", "http://localhost", "http://localhost:1"} // CHANGE THIS TO THE FRONTEND URL
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins, // CHANGE THIS - ALLOW ONLY FRONTEND URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	return c.Handler(router)
}
