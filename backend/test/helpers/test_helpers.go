// /backend/test/helpers/test_helpers.go
package helpers

import (
	"fmt"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

// SetupTestDB initializes the test database using the common InitTestDatabase() method
func SetupTestDB(t *testing.T) {
	// Initialize the validator (necessary for validation logic in tests)
	utils.InitValidator()

	// Set up the test-specific configuration
	config.AppConfig = config.Config{
		JWT: config.JWTConfig{
			Secret:     "testsecret",
			Expiration: 3600, // 1 hour expiration for JWTs
		},
		Server: config.ServerConfig{
			Port: 8080, // Set any port necessary for your tests
		},
		Environment: config.Environment{
			Environment: "test", // Ensure the environment is set to 'test'
		},
	}

	// Initialize the test database (using SQLite or any test-specific setup)
	err := database.InitTestDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize the test database: %v", err)
	}

	fmt.Println("Test database initialized")

	seedTestData(t)

	t.Cleanup(func() {
		TeardownTestDB(t)
	})
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T) {
	err := database.DB.Migrator().DropTable("users", "pages")
	if err != nil {
		t.Fatalf("Failed to drop tables: %v", err)
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		t.Fatalf("Failed to get the underlying SQL DB: %v", err)
	}

	err = sqlDB.Close()
	if err != nil {
		t.Fatalf("Failed to close the test database connection: %v", err)
	}

	fmt.Println("Test database connection closed")
}

// SeedTestData seeds initial data into the test database
func seedTestData(t *testing.T) {
	hashedPassword, _ := security.HashPassword("password123")
	testUser := models.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
	}
	err := database.DB.Create(&testUser).Error
	assert.NoError(t, err, "Failed to seed test data")

	pages := []models.Page{
		{Title: "Go Programming", Content: "A comprehensive guide to Go programming.", Language: "en", Url: "/go-programming"},
		{Title: "Python Programming", Content: "Learn Python with examples.", Language: "en", Url: "/python-programming"},
		{Title: "Danish Guide", Content: "Guide to Danish culture and language.", Language: "da", Url: "/danish-guide"},
	}
	err = database.DB.Create(&pages).Error
	assert.NoError(t, err, "Failed to seed pages")
}
