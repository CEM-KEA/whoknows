package config

import (
	"fmt"
	"strings"
)

// LoadEnv loads the .env file into the viper configuration
func LoadEnv(v ViperConfig) error {
	v.SetConfigFile(".env")
	err := v.ReadInConfig()

	if err != nil {
		return fmt.Errorf("error reading .env file - %s", err)
	}

	return nil
}

// LoadYaml loads the config.yaml file into the viper configuration
func LoadYaml(v ViperConfig) error {
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/config")
	err := v.ReadInConfig()

	if err != nil {
		return fmt.Errorf("error reading config.yaml file - %s", err)
	}

	return nil
}

// LoadConfig loads the configuration from the config file and environment variables
func LoadConfig(v ViperConfig) error {
	// Load the .env file
	err := LoadEnv(v)
	if err != nil {
		return err
	}

	// Load the config.yaml file
	err = LoadYaml(v)
	if err != nil {
		return err
	}

	// Enable viper to read environment variables
	v.AutomaticEnv()

	// Set the environment prefix
	v.SetEnvPrefix("API")

	// Replace the - and . in the environment variables with _
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Set the default values
	v.SetDefault("server.port", 8080)

	// Load the config into the AppConfig struct
	err = v.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling configuration - %s", err)
	}

	fmt.Printf("Configuration loaded successfully: %+v\n", AppConfig)

	return nil
}
