package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token                 string `json:"token"`
	RequirePasswordChange bool   `json:"require_password_change"`
}

func sanitizeUsername(username string) string {
	sanitizedUsername := strings.ReplaceAll(username, "\n", "")
	sanitizedUsername = strings.ReplaceAll(sanitizedUsername, "\r", "")
	return sanitizedUsername
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
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.LogError(err, "Failed to decode request body", nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.Validate(request); err != nil {
		sanitizedUsername := sanitizeUsername(request.Username)
		utils.LogError(err, "Request validation failed", logrus.Fields{
			"username": sanitizedUsername,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, valid, err := services.CheckUserPassword(database.DB, request.Password, request.Username)
	if err != nil || !valid {
		sanitizedUsername := sanitizeUsername(request.Username)
		utils.LogWarn("Invalid user credentials", logrus.Fields{
			"username": sanitizedUsername,
		})
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := security.GenerateJWT(user.ID, user.Username)
	if err != nil {
		sanitizedUsername := sanitizeUsername(request.Username)
		utils.LogError(err, "Failed to generate token", logrus.Fields{
			"username": sanitizedUsername,
		})
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	jwtModel := models.JWT{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(&jwtModel).Error; err != nil {
		utils.LogError(err, "Failed to save token to database", logrus.Fields{
			"username": user.Username,
		})
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	if err := services.UpdateLastLogin(database.DB, user); err != nil {
		utils.LogError(err, "Failed to update last login", logrus.Fields{
			"username": user.Username,
		})
		http.Error(w, "Failed to update last login", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token:                 token,
		RequirePasswordChange: user.UpdatedAt.Before(time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC)), // Check if user changed password after incident on 31/10/2024
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.LogError(err, "Failed to encode login response", logrus.Fields{
			"username": user.Username,
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	utils.LogInfo("User logged in successfully", logrus.Fields{
		"username": user.Username,
		"userID":   user.ID,
	})
}
