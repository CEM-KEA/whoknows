package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv loads the .env file using godotenv and populates AppConfig
func LoadEnv() error {
	envFilePath := os.Getenv("ENV_FILE_PATH")

	// If the ENV_FILE_PATH environment variable is not set, default to .env in the current working directory
	if envFilePath == "" {
		envFilePath = ".env" // Default to .env in the current working directory
	}

	// Load the .env file
	err := godotenv.Load(envFilePath)
	if err != nil {
		return fmt.Errorf("error reading .env file - %s", err)
	}

	// Populate AppConfig from environment variables
	err = loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("error loading configuration from env - %s", err)
	}

	return nil
}

// loadConfigFromEnv maps environment variables to AppConfig
func loadConfigFromEnv() error {
	// Server configuration
	port, err := strconv.Atoi(getEnv("API_SERVER_PORT", "8080"))

	if err != nil {
		return fmt.Errorf("invalid SERVER_PORT value")
	}

	AppConfig.Server.Port = port

	// Database configuration
	AppConfig.Database.FilePath = getEnv("API_DATABASE_FILE_PATH", "./internal/database/whoknows.db")
	AppConfig.Database.Migrate = getEnvAsBool("API_DATABASE_MIGRATE", false)
	AppConfig.Database.Seed = getEnvAsBool("API_DATABASE_SEED", false)
	AppConfig.Database.SeedFilePath = getEnv("API_DATABASE_SEED_FILE_PATH", "./internal/database/pages.json")

	// JWT configuration
	AppConfig.JWT.Secret = getEnv("API_JWT_SECRET", "")
	expiration, err := strconv.Atoi(getEnv("API_JWT_EXPIRATION", "3600"))

	if err != nil {
		return fmt.Errorf("invalid JWT_EXPIRATION value")
	}

	AppConfig.JWT.Expiration = expiration

	// Environment configuration
	AppConfig.Environment.Environment = getEnv("API_APP_ENVIRONMENT", "development")

	// Pagination configuration
	limit, err := strconv.Atoi(getEnv("API_PAGINATION_LIMIT", "10"))

	if err != nil {
		return fmt.Errorf("invalid PAGINATION_LIMIT value")
	}

	offset, err := strconv.Atoi(getEnv("API_PAGINATION_OFFSET", "0"))

	if err != nil {
		return fmt.Errorf("invalid PAGINATION_OFFSET value")
	}

	AppConfig.Pagination.Limit = limit
	AppConfig.Pagination.Offset = offset

	// Log configuration
	AppConfig.Log.Level = getEnv("API_LOG_LEVEL", "debug")
	AppConfig.Log.Format = getEnv("API_LOG_FORMAT", "text")

	return nil
}

// Helper function to get a string environment variable with a fallback value
func getEnv(key string, fallback string) string {
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
