package api

import (
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/api/handlers"
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


	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"}, // CHANGE THIS - ALLOW ONLY FRONTEND URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	return c.Handler(router)
}



