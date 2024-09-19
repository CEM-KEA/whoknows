package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/database"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents the registration request payload
// @Description Register a new user
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "User data"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {string} string "Validation error"
// @Failure 500 {string} string "Failed to create user"
// @Router /api/register [post]
// RegisterHandler handles the registration of a new user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	db := database.DB
	req := RegisterRequest{}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	// validate the request
	if err := utils.Validate(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create a new user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// Save the user to the database
	err = services.CreateUser(db, &user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}


