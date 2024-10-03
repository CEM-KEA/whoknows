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
	err := loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("error loading configuration: %s", err)
	}

	return nil
}

// loadConfigFromEnv maps environment variables to AppConfig
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
			return fmt.Errorf("error loading %s: %v", key, err)
		}
	}

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
