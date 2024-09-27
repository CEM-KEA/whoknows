package handlers

import (
	"net/http"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//Get the authorization header to get the jwt token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No Authorization header found"))
		return
	}

	//Sepreate the Bearer and the token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Authorization header format"))
		return
	}
	token := parts[1]

	//Revoke the jwt token
	err := security.RevokeJWT(database.DB, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to revoke token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
