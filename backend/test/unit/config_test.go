package unit_test

import (
	"os"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

// Test LoadEnv function (success scenario)
func TestLoadEnvSuccess(t *testing.T) {
	// Set the ENV_FILE_PATH to the actual .env file location
	os.Setenv("ENV_FILE_PATH", "../test.env")

	// Call the LoadEnv function
	err := config.LoadEnv()

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert AppConfig is populated correctly
	assert.Equal(t, 8080, config.AppConfig.Server.Port)
	assert.Equal(t, "./internal/database/whoknows.db", config.AppConfig.Database.FilePath)
	assert.Equal(t, "mysecret", config.AppConfig.JWT.Secret)
	assert.Equal(t, "test", config.AppConfig.Environment.Environment)
	assert.Equal(t, 10, config.AppConfig.Pagination.Limit)
	assert.Equal(t, 0, config.AppConfig.Pagination.Offset)
	assert.Equal(t, "debug", config.AppConfig.Log.Level)
	assert.Equal(t, "text", config.AppConfig.Log.Format)

	// Clean up environment variables to avoid test contamination
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DATABASE_FILE_PATH")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("APP_ENVIRONMENT")
	os.Unsetenv("PAGINATION_LIMIT")
	os.Unsetenv("PAGINATION_OFFSET")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_FORMAT")

	// Unset ENV_FILE_PATH to avoid test contamination
	os.Unsetenv("ENV_FILE_PATH")
}

// Test LoadEnv function (failure scenario)
func TestLoadEnvFailure(t *testing.T) {
	// Set the ENV_FILE_PATH to an invalid path to simulate failure
	os.Setenv("ENV_FILE_PATH", "/nonexistent.env")

	// Call LoadEnv
	err := config.LoadEnv()

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading .env file")

	// Unset ENV_FILE_PATH to avoid test contamination
	os.Unsetenv("ENV_FILE_PATH")
}
