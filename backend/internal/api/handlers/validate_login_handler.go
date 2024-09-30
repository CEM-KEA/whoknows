package handlers

import (
	"net/http"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
)

//	@Description	Validates the jwt token
//	@Security		Bearer
//	@Success		200	{string}	string	"valid"
//	@Failure		401	{string}	string	"No Authorization header found"
//	@Failure		401	{string}	string	"Invalid Authorization header format"
//	@Failure		401	{string}	string	"Invalid token"
//	@Failure		401	{string}	string	"Token expired/revoked"
//	@Router			/api/validate-login [get]
// Handler for validating the jwt token
func ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {

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

	tokenString := parts[1]

	_, err := security.ValidateJWT(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	err = security.ValidateJWTRevoked(database.DB, tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Token expired/revoked"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("valid"))
}
