package database

import (
	"fmt"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the variable that holds the database connection
var DB *gorm.DB

// InitDatabase initializes the Postgres database connection for production
func InitDatabase() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Name,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to Postgres database: %s", err)
	}

	fmt.Println("Postgres database connection established")
	return autoMigrate()
}

// InitTestDatabase initializes an SQLite in-memory database for testing
func InitTestDatabase() error {
	databasePath := config.AppConfig.TestDatabase.FilePath
	var err error
	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to SQLite test database: %s", err)
	}

	fmt.Println("SQLite test database connection established")
	return autoMigrate()
}

// autoMigrate migrates the database schema to the latest version
func autoMigrate() error {
	err := DB.AutoMigrate(&models.User{}, &models.Page{}, &models.JWT{})
	if err != nil {
		return fmt.Errorf("error migrating database: %s", err)
	}

	fmt.Println("Database schema migrated successfully")
	return nil
}
