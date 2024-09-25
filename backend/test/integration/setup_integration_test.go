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
    "github.com/CEM-KEA/whoknows/backend/internal/utils"
    "github.com/stretchr/testify/assert"
)

// setupTestDB initializes a fresh test database and seeds initial data
func setupTestDB(t *testing.T) {
    utils.InitValidator()

    config.AppConfig = config.Config{
        JWT: config.JWTConfig{
            Secret:     "testsecret",
            Expiration: 3600,
        },
        Server: config.ServerConfig{
            Port: 8080,
        },
        Database: config.DatabaseConfig{
            FilePath: ":memory:",
        },
        Environment: config.Environment{
            Environment: "test",
        },
    }

    err := database.InitDatabase()
    if err != nil {
        t.Fatalf("Failed to initialize the Test database: %v", err)
    }

    fmt.Println("Test Database initialized")

    seedTestData(t)

    t.Cleanup(func() {
        teardownTestDB(t)
    })
}

// seedTestData seeds the database with initial data for testing
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

// teardownTestDB closes and cleans up the test database
func teardownTestDB(t *testing.T) {
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
        t.Fatalf("Failed to close the Test database connection: %v", err)
    }

    fmt.Println("Test Database connection closed")
}

// TestMain sets up the test environment
func TestMain(m *testing.M) {
    os.Setenv("ENV_FILE_PATH", "../test.env")
    err := config.LoadEnv()
    if err != nil {
        log.Fatalf("Failed to load environment: %v", err)
    }

    code := m.Run()
    
    os.Exit(code)
}
