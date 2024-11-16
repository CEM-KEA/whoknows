package database

import (
	"fmt"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB holds the database connection
var DB *gorm.DB

// InitDatabase initializes the Postgres database connection for production
func InitDatabase() error {
	utils.LogInfo("Setting up database connection", nil)

	// Create the Data Source Name (DSN) for Postgres
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Name,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
	)

	// Open the Postgres database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.LogError(err, "Failed to connect to Postgres database", nil)
		return fmt.Errorf("error connecting to Postgres database: %s", err)
	}

	// Register the database callbacks for Prometheus metrics
	RegisterCallbacks(DB)

	utils.LogInfo("Database connection established", nil)

	// Perform migrations if enabled in the configuration
	if config.AppConfig.Database.Migrate {
		return autoMigrate()
	}

	utils.LogInfo("Database schema migration skipped", nil)
	return nil
}

// InitTestDatabase initializes an SQLite in-memory database for testing
func InitTestDatabase() error {
	utils.LogInfo("Setting up SQLite in-memory test database", nil)

	// Open the SQLite in-memory database connection
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		utils.LogError(err, "Failed to connect to SQLite in-memory database", nil)
		return fmt.Errorf("error connecting to SQLite in-memory database: %s", err)
	}

	utils.LogInfo("SQLite in-memory test database connection established", nil)

	// Perform migrations for the test database
	return autoMigrate()
}

// autoMigrate migrates the database schema to the latest version
func autoMigrate() error {
	utils.LogInfo("Migrating database schema", nil)

	// Define the migration
	m := gormigrate.New(DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: time.Now().Format("20060102150405"),
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{}, &models.Page{}, &models.JWT{}, &models.SearchLog{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.User{}, &models.Page{}, &models.JWT{}, &models.SearchLog{})
			},
		},
	})

	// Run the migration
	err := m.Migrate()
	if err != nil {
		utils.LogError(err, "Failed to migrate database schema", nil)
		return fmt.Errorf("error migrating database: %s", err)
	}

	utils.LogInfo("Database schema migration successful", nil)
	return nil
}
