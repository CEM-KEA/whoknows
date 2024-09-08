package database

import (
	"encoding/json"
	"fmt"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"gorm.io/gorm"
	"os"
	"time"
)

// PageData represents the structure of the JSON data
type PageData struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Language    string `json:"language"`
	LastUpdated string `json:"last_updated"`
	Content     string `json:"content"`
}

// Parse JSON data
func parseJSONData(file []byte) ([]PageData, error) {
	var pages []PageData
	err := json.Unmarshal(file, &pages)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %w", err)
	}

	return pages, nil
}

// SeedData seeds the database with initial data
func SeedData(db *gorm.DB) error {
	// Read JSON file
	file, err := os.ReadFile("./internal/database/pages.json")

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
		// Parse the last_updated field to time.Time
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
