package handlers

import (
	"net/http"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
)

//	@Description	Validates the jwt token
//	@Tags Authentication
//	@Security		Bearer
//	@Success		200	{string}	string	"valid"
//	@Failure		401	{string}	string	"No Authorization header found"
//	@Failure		401	{string}	string	"Invalid Authorization header format"
//	@Failure		401	{string}	string	"Invalid token"
//	@Failure		401	{string}	string	"Token expired/revoked"
//	@Router			/api/validate-login [get]
//
// Handler for validating the jwt token
func ValidateLoginHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Processing validate login request", nil)

	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.LogWarn("Authorization header is missing", nil)
		http.Error(w, "No Authorization header found", http.StatusUnauthorized)
		return
	}

	// Parse the Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.LogWarn("Invalid Authorization header format", logrus.Fields{
			"authHeader": authHeader,
		})
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	// Validate the JWT
	_, err := security.ValidateJWT(tokenString)
	if err != nil {
		utils.LogWarn("Invalid JWT token", logrus.Fields{
			"token": tokenString,
			"error": err.Error(),
		})
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if the token is revoked
	err = security.ValidateJWTRevoked(database.DB, tokenString)
	if err != nil {
		utils.LogWarn("JWT token is revoked or expired", logrus.Fields{
			"token": tokenString,
			"error": err.Error(),
		})
		http.Error(w, "Token expired/revoked", http.StatusUnauthorized)
		return
	}

	// Successful validation
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("valid"))

	utils.LogInfo("Token validation successful - user is logged in", logrus.Fields{
		"token": tokenString,
	})
}
