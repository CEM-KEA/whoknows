package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Status                string `json:"status"`
	Token                 string `json:"token"`
	RequirePasswordChange bool   `json:"require_password_change"`
}

// Login handles the login request.
//
//	@Summary Login a user
//	@Description Authenticate user and return a JWT token for further requests.
//	@Tags Authentication
//	@Accept json
//	@Produce json
//	@Param loginRequest body handlers.LoginRequest true "Login request body"
//	@Success 200 {object} handlers.LoginResponse "Successful login"
//	@Failure 400 {object} map[string]string "Invalid request body"
//	@Failure 401 {object} map[string]string "Invalid username or password"
//	@Failure 500 {object} map[string]string "Internal server error"
//	@Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Processing login request", nil)

	// Decode and sanitize the request body
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.LogError(err, "Failed to decode request body", nil)
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	utils.SanitizeStruct(&request)

	// Validate the sanitized request
	if err := utils.Validate(request); err != nil {
		utils.LogError(err, "Request validation failed", nil)
		utils.WriteJSONError(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// Check user credentials
	user, valid, err := services.CheckUserPassword(database.DB, request.Password, request.Username)
	if err != nil || !valid {
		utils.LogWarn("Invalid user credentials", nil)
		utils.WriteJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := security.GenerateJWT(user.ID, user.Username)
	if err != nil {
		utils.LogError(err, "Failed to generate token", nil)
		utils.WriteJSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Save token to database
	jwtModel := models.JWT{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(&jwtModel).Error; err != nil {
		utils.LogError(err, "Failed to save token to database", nil)
		utils.WriteJSONError(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	// Update last login timestamp
	if err := services.UpdateLastLogin(database.DB, user); err != nil {
		utils.LogError(err, "Failed to update last login", nil)
		utils.WriteJSONError(w, "Failed to update last login", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := LoginResponse{
		Token:                 token,
		RequirePasswordChange: user.UpdatedAt.Before(time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC)), // Check if user changed password after incident on 31/10/2024
	}

	// Send success response
	utils.JSONSuccess(w, map[string]interface{}{
		"status":                  "success",
		"token":                   response.Token,
		"require_password_change": response.RequirePasswordChange,
	}, http.StatusOK)

	utils.LogInfo("User logged in successfully", nil)
}
