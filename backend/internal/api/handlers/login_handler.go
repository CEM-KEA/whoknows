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
	Token string `json:"token"`
}

// LoginRequest represents the login request payload
// @Description Login with username and password
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Invalid username or password"
// @Router /api/login [post]
// Handler for login
func Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	//Decode the request body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Validate the request
	err = utils.Validate(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Retrieve the user from the database by username
	user, err := services.GetUserByUsername(database.DB, request.Username)

	if err != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	//Check if the user exists and if the password matches
	if !security.CheckPasswordHash(request.Password, user.PasswordHash) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	//Generate a JWT token for the user
	token, err := security.GenerateJWT(user.ID, user.Username)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	jwtModel := models.JWT{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		CreatedAt: time.Now(),
		RevokedAt: nil,
	}

	err = database.DB.Create(&jwtModel).Error
	if err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	//Return the token in the response
	response := LoginResponse{
		Token: token,
	}

	//Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
