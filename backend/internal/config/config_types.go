package config

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
