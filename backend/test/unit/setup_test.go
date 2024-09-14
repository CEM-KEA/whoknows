package unit_test

import (
	"fmt"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/database"
)

func setupTestDB(t *testing.T) {
	config.AppConfig = config.Config{
		JWT: config.JWTConfig{
			Secret:     "testsecret", // Ensure the same secret is used consistently
			Expiration: 3600,
		},
		Server: config.ServerConfig{
			Port: 8080,
		},
		Database: config.DatabaseConfig{
			// Use in-memory database
			FilePath: ":memory:",
		},
		Environment: config.Environment{
			Environment: "test",
		},
	}

	// Initialize the database
	err := database.InitDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize the Test database: %v", err)
	}

	fmt.Println("Test Database initialized")

	// Cleanup after the test completes
	t.Cleanup(func() {
		teardownTestDB(t)
	})
}

func teardownTestDB(t *testing.T) {
	// Get the underlying SQL DB connection from the GORM DB
	sqlDB, err := database.DB.DB()
	if err != nil {
		t.Fatalf("Failed to get the underlying SQL DB: %v", err)
	}

	// Close the SQL DB connection effectively ending the test and clearing the in-memory database
	err = sqlDB.Close()
	if err != nil {
		t.Fatalf("Failed to close the Test database connection: %v", err)
	}

	fmt.Println("Test Database connection closed")
}
