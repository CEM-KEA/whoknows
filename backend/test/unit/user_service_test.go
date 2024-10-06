package unit_test

import (
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	helpers.SetupTestDB(t)

	user := &models.User{
		Username: "test_user",
		Email:    "test_user@example.com",
	}

	err := services.CreateUser(database.DB, user)
	assert.NoError(t, err)

	var result models.User
	database.DB.Where("username = ?", user.Username).First(&result)

	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Email, result.Email)
}

func TestGetUserByUsername(t *testing.T) {
	helpers.SetupTestDB(t)

	username := "test_user"
	expectedUser := &models.User{
		Username: "test_user",
		Email:    "test_user@example.com",
	}
	database.DB.Create(expectedUser)

	result, err := services.GetUserByUsername(database.DB, username)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Email, result.Email)
}
