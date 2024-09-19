package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Handler for login
func Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	//Decode the request body into the LoginRequest struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Retrieve the user from the database by email
	user, err := services.GetUserByEmail(database.DB, request.Email)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	//Check if the user exists and if the password matches
	if !security.CheckPasswordHash(request.Password, user.PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	//Generate a JWT token for the user
	token, err := security.GenerateJWT(user.ID, user.Email)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
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
