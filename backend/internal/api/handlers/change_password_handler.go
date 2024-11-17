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

// ChangePasswordRequest represents the change password request payload
//
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

	var request ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.LogError(err, "Failed to decode request body", nil)
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	utils.SanitizeStruct(&request)

	if err := utils.Validate(request); err != nil {
		utils.LogError(err, "Request validation failed", nil)
		utils.WriteJSONError(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	_, valid, err := services.CheckUserPassword(database.DB, request.Password, request.Username)
	if err != nil || !valid {
		utils.LogWarn("Invalid user credentials", nil)
		utils.WriteJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError(err, "Password hashing failed", nil)
		utils.WriteJSONError(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	user, _ := services.GetUserByUsername(database.DB, request.Username)
	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now()
	if err := services.UpdateUser(database.DB, user); err != nil {
		utils.LogError(err, "Failed to update user in database", nil)
		utils.WriteJSONError(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("Password changed successfully", nil)
	utils.JSONSuccess(w, map[string]interface{}{
		"status":  "success",
		"message": "Password changed successfully",
	}, http.StatusOK)
}
