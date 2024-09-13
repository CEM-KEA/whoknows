package unit_test

import (
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Set up the test database
	setupTestDB(t)

	// Create a new user object
	user := &models.User{Email: "test@example.com"}

	// Call the CreateUser service to create the user
	err := services.CreateUser(database.DB, user)

	// Validate that there were no errors
	assert.NoError(t, err)

	// Fetch the user from the database to verify it was created
	var result models.User

	database.DB.Where("email = ?", user.Email).First(&result)

	// Validate the result
	assert.Equal(t, user.Email, result.Email)
}

func TestGetUserByEmail(t *testing.T) {
	// Set up the test database
	setupTestDB(t)

	// Insert a user into the in-memory database
	email := "test@example.com"
	expectedUser := &models.User{Email: email}
	database.DB.Create(expectedUser)

	// Call the actual service method
	result, err := services.GetUserByEmail(database.DB, email)

	// Validate the results
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, result.Email)
}
