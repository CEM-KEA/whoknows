package unit_test

import (
	"os"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	. "github.com/CEM-KEA/whoknows/backend/internal/database"
)

func TestInitDatabase(t *testing.T) {
	// Setup a temporary database file
	tempFile, err := os.CreateTemp("", "testdb_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	// Set the config to use the temporary database file
	config.AppConfig.Database.FilePath = tempFile.Name()

	// Initialize the database
	err = InitDatabase()
	if err != nil {
		t.Fatalf("InitDatabase failed: %s", err)
	}

	// Check if the DB variable is not nil
	if DB == nil {
		t.Fatal("DB is nil after InitDatabase")
	}

	// Check if the connection is valid
	sqlDB, err := DB.DB()
	if err != nil {
		t.Fatalf("Failed to get sqlDB from gorm DB: %s", err)
	}
	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %s", err)
	}
}

func TestMigrateDatabase(t *testing.T) {
	// Setup a temporary database file
	tempFile, err := os.CreateTemp("", "testdb_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	// Set the config to use the temporary database file
	config.AppConfig.Database.FilePath = tempFile.Name()

	// Initialize the database
	err = InitDatabase()
	if err != nil {
		t.Fatalf("InitDatabase failed: %s", err)
	}

	// Migrate the database
	err = MigrateDatabase()
	if err != nil {
		t.Fatalf("MigrateDatabase failed: %s", err)
	}

	// Check if the tables were created
	tables := []string{"users", "pages"}
	for _, table := range tables {
		if !DB.Migrator().HasTable(table) {
			t.Fatalf("Table %s does not exist after migration", table)
		}
	}
}
