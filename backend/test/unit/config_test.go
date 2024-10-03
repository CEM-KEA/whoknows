package unit_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestLoadEnvSuccess tests if LoadEnv successfully loads variables from a .env file
func TestLoadEnvSuccess(t *testing.T) {
	os.Setenv("ENV_FILE_PATH", "../test.env")

	err := config.LoadEnv()
	assert.NoError(t, err)
	assert.Equal(t, 8080, config.AppConfig.Server.Port)
	assert.Equal(t, "localhost", config.AppConfig.Database.Host)
	assert.Equal(t, 5432, config.AppConfig.Database.Port)
	assert.Equal(t, "user", config.AppConfig.Database.User)
	assert.Equal(t, "password", config.AppConfig.Database.Password)
	assert.Equal(t, "mydb", config.AppConfig.Database.Name)
	assert.Equal(t, "disable", config.AppConfig.Database.SSLMode)
	assert.Equal(t, true, config.AppConfig.Database.Migrate)
	assert.Equal(t, "mysecret", config.AppConfig.JWT.Secret)
	assert.Equal(t, 3600, config.AppConfig.JWT.Expiration)
	assert.Equal(t, "test", config.AppConfig.Environment.Environment)
	assert.Equal(t, 10, config.AppConfig.Pagination.Limit)
	assert.Equal(t, 0, config.AppConfig.Pagination.Offset)
	assert.Equal(t, "debug", config.AppConfig.Log.Level)
	assert.Equal(t, "text", config.AppConfig.Log.Format)
	assert.Equal(t, "weatherapikey", config.AppConfig.WeatherAPI.OpenWeatherAPIKey)

	os.Unsetenv("ENV_FILE_PATH")
}

// TestLoadEnvFailure simulates a failure when required environment variables are missing
func TestLoadEnvFailure(t *testing.T) {
	// Clear all environment variables to simulate a missing .env and no environment variables
	os.Clearenv()

	// Optionally, print environment variables to confirm they are cleared
	fmt.Println("Environment after clearing:")
	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	// Call LoadEnv and expect it to fail due to missing required environment variables
	err := config.LoadEnv()

	// Assert that an error occurred
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "error loading configuration")
}
