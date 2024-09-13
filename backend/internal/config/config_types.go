package config

// Config is the struct that holds the application configuration
type Config struct {
	Environment Environment
	JWT         JWTConfig
	Server      ServerConfig
	Database    DatabaseConfig
}

// Environment is the struct that holds the environment configuration
type Environment struct {
	Environment string
}

// JWTConfig is the struct that holds the JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration int
}

// ServerConfig is the struct that holds the server configuration
type ServerConfig struct {
	Port int
}

// DatabaseConfig is the struct that holds the database configuration for sqlite3
type DatabaseConfig struct {
	FilePath     string
	Migrate      bool
	Seed         bool
	SeedFilePath string
}

// AppConfig is the variable that holds the application configuration
var AppConfig Config
