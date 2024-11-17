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


// DB is a global variable that holds the database connection instance
// managed by GORM. It is used throughout the application to interact
// with the database.
var DB *gorm.DB


// InitDatabase initializes the database connection using the configuration
// specified in the application configuration. It sets up the Data Source Name (DSN)
// for Postgres, opens the database connection using GORM, registers database
// callbacks for Prometheus metrics, and performs schema migrations if enabled
// in the configuration.
//
// Returns an error if the database connection fails or if there is an error
// during schema migration.
func InitDatabase() error {
	utils.LogInfo("Setting up database connection", nil)
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
		utils.LogError(err, "Failed to connect to Postgres database", nil)
		return fmt.Errorf("error connecting to Postgres database: %s", err)
	}

	// Register the database callbacks for Prometheus metrics
	RegisterCallbacks(DB)

	utils.LogInfo("Database connection established", nil)

	if config.AppConfig.Database.Migrate {
		return autoMigrate()
	}

	utils.LogInfo("Database schema migration skipped", nil)
	return nil
}

// InitTestDatabase initializes an in-memory SQLite database for testing purposes.
// It sets up the database connection using GORM and performs auto-migration.
//
// Returns an error if the database connection or migration fails.
func InitTestDatabase() error {
	utils.LogInfo("Setting up SQLite in-memory test database", nil)
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		utils.LogError(err, "Failed to connect to SQLite in-memory database", nil)
		return fmt.Errorf("error connecting to SQLite in-memory database: %s", err)
	}

	utils.LogInfo("SQLite in-memory test database connection established", nil)

	return autoMigrate()
}

// autoMigrate handles the automatic migration of the database schema.
// It logs the start and end of the migration process, and performs the migration
// using gormigrate. If the migration fails, it logs the error and returns it.
//
// Returns an error if the migration fails.
func autoMigrate() error {
	utils.LogInfo("Migrating database schema", nil)
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

	err := m.Migrate()
	if err != nil {
		utils.LogError(err, "Failed to migrate database schema", nil)
		return fmt.Errorf("error migrating database: %s", err)
	}

	utils.LogInfo("Database schema migration successful", nil)
	return nil
}
