package handlers

import (
	"net/http"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
)

// LogoutHandler logs out the user by revoking the jwt token
//
//	@Description	Logs out the user by revoking the jwt token
//	@Tags Authentication
//	@Security		Bearer
//	@Success		200	{string}	string	"Logged out successfully"
//	@Failure		401	{string}	string	"No Authorization header found"
//	@Failure		401	{string}	string	"Invalid Authorization header format"
//	@Failure		500	{string}	string	"Failed to revoke token"
//	@Router			/api/logout [get]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Processing logout request", nil)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.LogWarn("No Authorization header found", nil)
		utils.WriteJSONError(w, "No Authorization header found", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.LogWarn("Invalid Authorization header format", nil)
		utils.WriteJSONError(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}
	token := parts[1]

	err := security.RevokeJWT(database.DB, token)
	if err != nil {
		utils.LogError(err, "Failed to revoke token", nil)
		utils.WriteJSONError(w, "Failed to revoke token", http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, map[string]interface{}{
		"status":  "success",
		"message": "Logged out successfully",
	}, http.StatusOK)

	utils.LogInfo("User logged out successfully", nil)
}
