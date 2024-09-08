package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// Config is the struct that holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig is the struct that holds the server configuration
type ServerConfig struct {
	Port int
}

// DatabaseConfig is the struct that holds the database configuration for sqlite3
type DatabaseConfig struct {
	FilePath string
	Migrate  bool
	Seed     bool
}

// AppConfig is the variable that holds the application configuration
var AppConfig Config

// LoadEnv loads the .env file into the environment variables
func LoadEnv() error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		return fmt.Errorf("error reading .env file - %s", err)
	}

	return nil
}

func LoadYaml() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	err := viper.ReadInConfig()

	if err != nil {
		return fmt.Errorf("error reading config.yaml file - %s", err)
	}

	return nil
}

// LoadConfig loads the configuration from the config file
func LoadConfig() error {
	// load the .env file
	err := LoadEnv()

	if err != nil {
		return err
	}

	// load the config.yaml file
	err = LoadYaml()

	if err != nil {
		return err
	}

	// enable viper to read environment variables
	viper.AutomaticEnv()

	// set the environment prefix
	viper.SetEnvPrefix("API")

	// replace the - and . in the environment variables with _
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// set the default values
	viper.SetDefault("server.port", 8080)

	// load the config into the AppConfig struct
	err = viper.Unmarshal(&AppConfig)

	if err != nil {
		return fmt.Errorf("error unmarshalling configuration - %s", err)
	}

	fmt.Printf("Configuration loaded successfully: %+v\n", AppConfig)

	return nil
}
