package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordRequest struct {
	Username          string `json:"username" validate:"required"`
	Password          string `json:"old_password" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

//ChangePasswordRequest represents the change password request payload
//	@Summary Change user password
//	@Description Endpoint to change the password of a user
//	@Tags Authentication
//	@Accept json
//	@Produce json
//	@Param changePasswordRequest body handlers.ChangePasswordRequest true "Change password payload"
//	@Success 200 {object} map[string]string "Password changed successfully"
//	@Failure 400 {object} map[string]string "Validation error"
//	@Failure 500 {object} map[string]string "Failed to change password"
//	@Router /api/change-password [post]
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Processing change password request", nil)

	// Decode request body
	var request ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.LogError(err, "Failed to decode request body", nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request data
	if err := utils.Validate(request); err != nil {
		utils.LogError(err, "Request validation failed", nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check user credentials
	user, valid, err := services.CheckUserPassword(database.DB, request.Password, request.Username)
	if err != nil || !valid {
		utils.LogWarn("Invalid user credentials", map[string]interface{}{
			"username": request.Username,
		})
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError(err, "Password hashing failed", map[string]interface{}{
			"username": request.Username,
		})
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	// Update user password and timestamp
	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now()
	if err := services.UpdateUser(database.DB, user); err != nil {
		utils.LogError(err, "Failed to update user in database", map[string]interface{}{
			"username": request.Username,
		})
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	// Respond with success
	utils.LogInfo("Password changed successfully", map[string]interface{}{
		"username": request.Username,
	})
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"}); err != nil {
		utils.LogError(err, "Failed to encode response", nil)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
