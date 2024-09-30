package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv loads the configuration based on the environment
func LoadEnv() error {
	envFilePath := os.Getenv("ENV_FILE_PATH")

	// Load .env file only if ENV_FILE_PATH is set (for local dev)
	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			return fmt.Errorf("error reading .env file - %s", err)
		}
		fmt.Println("Loaded configuration from .env file")
	} else {
		fmt.Println("Loading configuration from environment variables")
	}

	// Populate AppConfig from environment variables
	return loadConfigFromEnv()
}

// loadConfigFromEnv maps environment variables to AppConfig
func loadConfigFromEnv() error {
	// General Config
	port, err := strconv.Atoi(getEnv("API_SERVER_PORT", "8080"))
	if err != nil {
		return fmt.Errorf("invalid API_SERVER_PORT value")
	}
	AppConfig.Server.Port = port

	// Database Configuration
	AppConfig.Database.Host = getEnv("API_DATABASE_HOST", "localhost")
	AppConfig.Database.Port, err = strconv.Atoi(getEnv("API_DATABASE_PORT", "5432"))
	if err != nil {
		return fmt.Errorf("invalid API_DATABASE_PORT value")
	}
	AppConfig.Database.User = getEnv("API_DATABASE_USER", "postgres")
	AppConfig.Database.Password = getEnv("API_DATABASE_PASSWORD", "")
	AppConfig.Database.Name = getEnv("API_DATABASE_NAME", "whoknows")
	AppConfig.Database.SSLMode = getEnv("API_DATABASE_SSL_MODE", "disable")
	AppConfig.Database.Migrate = getEnvAsBool("API_DATABASE_MIGRATE", true)
	AppConfig.Database.Seed = getEnvAsBool("API_DATABASE_SEED", false) // Default to false for production

	// Test Database Configuration
	AppConfig.TestDatabase.FilePath = getEnv("API_TEST_DATABASE_FILE_PATH", ":memory:")

	// JWT Configuration
	AppConfig.JWT.Secret = getEnv("API_JWT_SECRET", "mysecret")
	AppConfig.JWT.Expiration, err = strconv.Atoi(getEnv("API_JWT_EXPIRATION", "3600"))
	if err != nil {
		return fmt.Errorf("invalid API_JWT_EXPIRATION value")
	}

	// Environment Configuration
	AppConfig.Environment.Environment = getEnv("API_ENVIRONMENT", "development")

	// Pagination Configuration
	AppConfig.Pagination.Limit, err = strconv.Atoi(getEnv("API_PAGINATION_LIMIT", "10"))
	if err != nil {
		return fmt.Errorf("invalid API_PAGINATION_LIMIT value")
	}
	AppConfig.Pagination.Offset, err = strconv.Atoi(getEnv("API_PAGINATION_OFFSET", "0"))
	if err != nil {
		return fmt.Errorf("invalid API_PAGINATION_OFFSET value")
	}

	// Log Configuration
	AppConfig.Log.Level = getEnv("API_LOG_LEVEL", "debug")
	AppConfig.Log.Format = getEnv("API_LOG_FORMAT", "text")

	// Weather API Configuration
	AppConfig.WeatherAPI.OpenWeatherAPIKey = getEnv("API_WEATHER_API_KEY", "")
	if AppConfig.WeatherAPI.OpenWeatherAPIKey == "" {
		return fmt.Errorf("API_WEATHER_API_KEY is required")
	}

	return nil
}

// Helper function to get a string environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to get a boolean environment variable
func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1"
	}
	return fallback
}
