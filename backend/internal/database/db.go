package database

import (
	"fmt"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the variable that holds the database connection
var DB *gorm.DB

// InitDatabase creates a new database file if it does not exist and establishes a connection to it
func InitDatabase() error {
	databasePath := config.AppConfig.Database.FilePath

	fmt.Printf("Connection to Database at: %s\n", databasePath)

	var err error
	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("error connecting to database: %s", err)
	}

	fmt.Println("Database connection established")

	return nil
}

// MigrateDatabase migrates the database schema to the latest version
func MigrateDatabase() error {
	err := DB.AutoMigrate(&models.User{}, &models.Page{})

	if err != nil {
		return fmt.Errorf("error migrating database: %s", err)
	}

	fmt.Println("Database migrated successfully")

	return nil
}
