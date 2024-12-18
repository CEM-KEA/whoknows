package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// LoadEnv loads the environment configuration for the application.
// It first checks if an environment file path is specified in the ENV_FILE_PATH
// environment variable. If specified, it attempts to load the configuration
// from the .env file at that path. If the .env file cannot be read, an error
// is logged and returned.
//
// If no .env file path is specified, it loads the configuration directly from
// the environment variables.
//
// The function also populates the AppConfig from the environment variables.
// If there is an error during this process, it logs the error and returns it.
//
// Returns an error if there is any issue loading the configuration.
func LoadEnv() error {
	utils.LogInfo("Loading environment configuration", nil)

	envFilePath := os.Getenv("ENV_FILE_PATH")

	// Load .env file if specified
	if envFilePath != "" {
		if err := godotenv.Load(envFilePath); err != nil {
			utils.LogError(err, "Error reading .env file", logrus.Fields{
				"path": envFilePath,
			})
			return fmt.Errorf("error reading .env file: %w", err)
		}
		utils.LogInfo("Loaded configuration from .env file", logrus.Fields{
			"path": envFilePath,
		})
	} else {
		utils.LogInfo("Loading configuration from environment variables", nil)
	}

	// Populate AppConfig from environment variables
	if err := loadConfigFromEnv(); err != nil {
		utils.LogError(err, "Error loading configuration from environment variables", nil)
		return fmt.Errorf("error loading configuration: %w", err)
	}

	utils.LogInfo("Configuration loaded successfully", nil)
	return nil
}


// loadConfigFromEnv loads the application configuration from environment variables.
// It maps environment variable keys to their respective configuration fields and
// attempts to load each one using helper functions such as getEnv, getEnvAsInt, and getEnvAsBool.
// If any environment variable fails to load, it logs the error and returns it.
// If all environment variables are loaded successfully, it logs a success message.
//
// Returns an error if any environment variable fails to load.
func loadConfigFromEnv() error {
	var err error

	loadConfig := map[string]func() error{
		// General Config
		"API_SERVER_PORT": func() error { AppConfig.Server.Port, err = getEnvAsInt("API_SERVER_PORT"); return err },

		// Database Configuration
		"API_DATABASE_HOST":     func() error { AppConfig.Database.Host, err = getEnv("API_DATABASE_HOST"); return err },
		"API_DATABASE_PORT":     func() error { AppConfig.Database.Port, err = getEnvAsInt("API_DATABASE_PORT"); return err },
		"API_DATABASE_USER":     func() error { AppConfig.Database.User, err = getEnv("API_DATABASE_USER"); return err },
		"API_DATABASE_PASSWORD": func() error { AppConfig.Database.Password, err = getEnv("API_DATABASE_PASSWORD"); return err },
		"API_DATABASE_NAME":     func() error { AppConfig.Database.Name, err = getEnv("API_DATABASE_NAME"); return err },
		"API_DATABASE_SSL_MODE": func() error { AppConfig.Database.SSLMode, err = getEnv("API_DATABASE_SSL_MODE"); return err },
		"API_DATABASE_MIGRATE":  func() error { AppConfig.Database.Migrate, err = getEnvAsBool("API_DATABASE_MIGRATE"); return err },

		// JWT Configuration
		"API_JWT_SECRET":     func() error { AppConfig.JWT.Secret, err = getEnv("API_JWT_SECRET"); return err },
		"API_JWT_EXPIRATION": func() error { AppConfig.JWT.Expiration, err = getEnvAsInt("API_JWT_EXPIRATION"); return err },

		// Environment Configuration
		"API_ENVIRONMENT": func() error { AppConfig.Environment.Environment, err = getEnv("API_ENVIRONMENT"); return err },

		// Pagination Configuration
		"API_PAGINATION_LIMIT":  func() error { AppConfig.Pagination.Limit, err = getEnvAsInt("API_PAGINATION_LIMIT"); return err },
		"API_PAGINATION_OFFSET": func() error { AppConfig.Pagination.Offset, err = getEnvAsInt("API_PAGINATION_OFFSET"); return err },

		// Log Configuration
		"API_LOG_LEVEL":  func() error { AppConfig.Log.Level, err = getEnv("API_LOG_LEVEL"); return err },
		"API_LOG_FORMAT": func() error { AppConfig.Log.Format, err = getEnv("API_LOG_FORMAT"); return err },

		// Weather API Configuration
		"API_WEATHER_API_KEY": func() error { AppConfig.WeatherAPI.OpenWeatherAPIKey, err = getEnv("API_WEATHER_API_KEY"); return err },
	}

	for key, fn := range loadConfig {
		if err := fn(); err != nil {
			utils.LogError(err, "Error loading configuration", logrus.Fields{"key": key})
			return fmt.Errorf("error loading %s: %v", key, err)
		}
	}
	
	utils.LogInfo("All environment variables loaded successfully", nil)
	return nil
}

// Helper function to get a string environment variable
func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("%s environment variable is required", key)
}

// Helper function to get an integer environment variable
func getEnvAsInt(key string) (int, error) {
	valueStr, err := getEnv(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value", key)
	}
	return value, nil
}

// Helper function to get a boolean environment variable
func getEnvAsBool(key string) (bool, error) {
	valueStr, err := getEnv(key)
	if err != nil {
		return false, err
	}
	return valueStr == "true" || valueStr == "1", nil
}
