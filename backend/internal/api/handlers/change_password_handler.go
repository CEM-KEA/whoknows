package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordRequest struct {
	Username          string `json:"username" validate:"required"`
	Password          string `json:"old_password" validate:"required"`
	NewPassword  	    string `json:"new_password" validate:"required"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

// ChangePasswordRequest represents the change password request payload
//	@Description	Change password
//	@Accept			json
//	@Security		Bearer
//	@Produce		json
//	@Param			register	body		RegisterRequest	true	"User data"
//	@Success		200			{string}	string			"Password changed successfully"
//	@Failure		400			{string}	string			"Validation error"
//	@Failure		500			{string}	string			"Failed to change password"
//	@Router			/api/change-password [post]
// ChangePasswordHandler handles changing of a user password, requires the user to be logged in as the request uses the Bearer token to get user information
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request ChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = utils.Validate(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := services.GetUserByUsername(database.DB, request.Username)

	if err != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	if !security.CheckPasswordHash(request.Password, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now()

	err = services.UpdateUser(database.DB, user)
	if err != nil {
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Password changed successfully"))
}
