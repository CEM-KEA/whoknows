package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	// Password2 is used to confirm the password, it is optional, so it is omitted if it is not provided or an empty string
	// If it is provided, it must be equal to the Password field
	Password2 string `json:"password2" validate:"omitempty,eqfield=Password"`
}

// RegisterRequest represents the registration request payload
//
//	@Description	Register a new user
//	@Tags Authentication
//	@Accept			json
//	@Produce		json
//	@Param			register	body		RegisterRequest	true	"User data"
//	@Success		201			{string}	string			"User created successfully"
//	@Failure		400			{string}	string			"Validation error"
//	@Failure		500			{string}	string			"Failed to create user"
//	@Router			/api/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Processing register request", nil)
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError(err, "Failed to decode request body", nil)
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	utils.SanitizeStruct(&req)

	if err := utils.Validate(req); err != nil {
		utils.LogWarn("Request validation failed", logrus.Fields{"error": err.Error()})

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s: %s", err.Field(), err.Tag()))
		}
		validationMessage := strings.Join(validationErrors, "; ")
	
		http.Error(w, validationMessage, http.StatusBadRequest)
		return
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		utils.LogError(err, "Failed to hash password", nil)
		utils.WriteJSONError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := services.CreateUser(database.DB, &user); err != nil {
		utils.LogError(err, "Failed to create user in database", nil)
		utils.WriteJSONError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	utils.IncrementUserRegistrations()

	utils.JSONSuccess(w, map[string]interface{}{
		"status":  "success",
		"message": "User created successfully",
	}, http.StatusCreated)

	utils.LogInfo("User registered successfully", nil)
}
