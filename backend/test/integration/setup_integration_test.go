package integration_test

import (
    "fmt"
    "log"
    "os"
    "testing"

    "github.com/CEM-KEA/whoknows/backend/internal/config"
    "github.com/CEM-KEA/whoknows/backend/internal/database"
    "github.com/CEM-KEA/whoknows/backend/internal/models"
    "github.com/CEM-KEA/whoknows/backend/internal/security"
    "github.com/CEM-KEA/whoknows/backend/internal/utils" // Include utils for validator
    "github.com/stretchr/testify/assert"
)

// setupTestDB initializes a fresh test database and seeds initial data
func setupTestDB(t *testing.T) {
    // Initialize Validator
    utils.InitValidator() // Initialize the validator

    config.AppConfig = config.Config{
        JWT: config.JWTConfig{
            Secret:     "testsecret",
            Expiration: 3600,
        },
        Server: config.ServerConfig{
            Port: 8080,
        },
        Database: config.DatabaseConfig{
            FilePath: ":memory:", // Use in-memory database for integration tests
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

    // Seed the database with initial data if necessary
    seedTestData(t)

    // Cleanup after the test completes
    t.Cleanup(func() {
        teardownTestDB(t)
    })
}

// seedTestData seeds the database with initial data for testing
func seedTestData(t *testing.T) {
    // Seed the database with initial data for tests
    // For example, creating a test user
    hashedPassword, _ := security.HashPassword("password123")
    testUser := models.User{
        Username:     "testuser",
        PasswordHash: hashedPassword,
    }
    err := database.DB.Create(&testUser).Error
    assert.NoError(t, err, "Failed to seed test data")

    // Seed pages with unique URLs
    pages := []models.Page{
        {Title: "Go Programming", Content: "A comprehensive guide to Go programming.", Language: "en", Url: "/go-programming"},
        {Title: "Python Programming", Content: "Learn Python with examples.", Language: "en", Url: "/python-programming"},
        {Title: "Danish Guide", Content: "Guide to Danish culture and language.", Language: "da", Url: "/danish-guide"},
    }
    err = database.DB.Create(&pages).Error
    assert.NoError(t, err, "Failed to seed pages")
}

// teardownTestDB closes and cleans up the test database
func teardownTestDB(t *testing.T) {
    // Drop all tables to clean up the database
    err := database.DB.Migrator().DropTable("users", "pages")
    if err != nil {
        t.Fatalf("Failed to drop tables: %v", err)
    }

    // Get the underlying SQL DB connection from the GORM DB
    sqlDB, err := database.DB.DB()
    if err != nil {
        t.Fatalf("Failed to get the underlying SQL DB: %v", err)
    }

    // Close the SQL DB connection, effectively ending the test and clearing the in-memory database
    err = sqlDB.Close()
    if err != nil {
        t.Fatalf("Failed to close the Test database connection: %v", err)
    }

    fmt.Println("Test Database connection closed")
}

// TestMain sets up the test environment
func TestMain(m *testing.M) {
    // Load the test environment
    os.Setenv("ENV_FILE_PATH", "../test.env")
    err := config.LoadEnv()
    if err != nil {
        log.Fatalf("Failed to load environment: %v", err)
    }

    // Run the tests
    code := m.Run()

    // Exit with the appropriate code
    os.Exit(code)
}
