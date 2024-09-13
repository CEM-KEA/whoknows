package database

import (
	"encoding/json"
	"fmt"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"gorm.io/gorm"
	"os"
	"time"
)

// Page represents the structure of the page data in the JSON file
type Page struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Language    string `json:"language"`
	Content     string `json:"content"`
	LastUpdated string `json:"last_updated"`
}

// parseJSONData parses the JSON data from the file and returns a slice of Page
func parseJSONData(data []byte) ([]Page, error) {
	var pages []Page
	err := json.Unmarshal(data, &pages)

	if err != nil {
		return nil, err
	}

	return pages, nil
}

// SeedData seeds the database with initial data from the specified file
func SeedData(db *gorm.DB, filePath string) error {
	// Read JSON file
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}

	// Parse JSON data
	pages, err := parseJSONData(file)
	if err != nil {
		return fmt.Errorf("error parsing JSON data: %w", err)
	}

	users := []models.User{
		{Username: "admin", Email: "keamonk1@stud.kea.dk", PasswordHash: "$2a$12$aUm2cOPvlw.9u0LuIkmfy.D/3UajkIFSNx27dIH9IsQUzNJR/jv/a"},
	}

	// Seed the users
	for _, user := range users {
		err := db.Create(&user).Error
		if err != nil {
			return fmt.Errorf("error seeding user data: %w", err)
		}
	}

	// Seed the pages
	for _, page := range pages {
		lastUpdated, err := time.Parse("2006-01-02 15:04:05", page.LastUpdated)
		if err != nil {
			return fmt.Errorf("error parsing last_updated field for page %s: %w", page.Title, err)
		}

		err = db.Create(&models.Page{
			Title:     page.Title,
			Url:       page.Url,
			Language:  page.Language,
			Content:   page.Content,
			CreatedAt: time.Now(),
			UpdatedAt: lastUpdated,
		}).Error
		if err != nil {
			return fmt.Errorf("error seeding page data for page %s: %w", page.Title, err)
		}
	}

	fmt.Println("Data seeded successfully")

	return nil
}
