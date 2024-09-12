package middlewares

import (
	"net/http"
	"github.com/rs/cors"
)

// CORSMiddleware applies the CORS settings to the router
func CORSMiddleware(next http.Handler) http.Handler {
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // CHANGE THIS - ALLOW ONLY FRONTEND URL
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })
    return corsHandler.Handler(next)
}
