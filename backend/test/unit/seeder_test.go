package unit_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSeedData(t *testing.T) {
	// Set up the test database
	setupTestDB(t)

	// Debug print: Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get the current working directory:", err)
	}

	fmt.Println("Debug: Current working directory:", wd)

	// Construct the absolute path to test_pages.json
	absPath, err := filepath.Abs("../../internal/database/test_pages.json")
	assert.NoError(t, err)

	// Clear the users and pages table before starting the test
	database.DB.Exec("DELETE FROM users")
	database.DB.Exec("DELETE FROM pages")

	// Call the SeedData function with the absolute file path
	err = database.SeedData(database.DB, absPath)
	assert.NoError(t, err)

	// Debug print: Select all data from the `users` table
	var users []models.User

	database.DB.Find(&users)
	fmt.Println("Debug: Users in the database:")

	for _, user := range users {
		fmt.Printf("User: ID=%d, Username=%s, Email=%s\n", user.ID, user.Username, user.Email)
	}

	// Debug print: Select all data from the `pages` table
	var pages []models.Page

	database.DB.Find(&pages)
	fmt.Println("Debug: Pages in the database:")

	for _, page := range pages {
		fmt.Printf("Page: ID=%d, Title=%s, URL=%s, Language=%s\n", page.ID, page.Title, page.Url, page.Language)
	}

	// Validate that the user was seeded
	var user models.User
	err = database.DB.Where("email = ?", "keamonk1@stud.kea.dk").First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, "keamonk1@stud.kea.dk", user.Email)

	// Validate that the page was seeded
	var page models.Page
	err = database.DB.Where("title = ?", "Test Page").First(&page).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Page", page.Title)
	assert.Equal(t, "/test-page", page.Url)
	assert.Equal(t, "en", page.Language)
	assert.Equal(t, "This is a test page.", page.Content)
}
