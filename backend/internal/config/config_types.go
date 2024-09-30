package config

// Config is the struct that holds the application configuration
type Config struct {
	Environment  Environment
	JWT          JWTConfig
	Server       ServerConfig
	Database     DatabaseConfig
	TestDatabase TestDatabaseConfig
	Pagination   PaginationConfig
	Log          LogConfig
	WeatherAPI   WeatherAPIConfig
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
type TestDatabaseConfig struct {
	FilePath     string
}

// DatabaseConfig is the struct that holds the database configuration for Postgres
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	Migrate  bool
	Seed	 bool
}

// PaginationConfig holds the pagination-related configuration
type PaginationConfig struct {
	Limit  int
	Offset int
}

// LogConfig holds the logging configuration
type LogConfig struct {
	Level  string
	Format string
}

type WeatherAPIConfig struct {
	OpenWeatherAPIKey string
}

// AppConfig is the variable that holds the application configuration
var AppConfig Config
