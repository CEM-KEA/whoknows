package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/api"
	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
)

// @title						WhoKnows API
// @version					1.0
// @description				This is the API for the WhoKnows application
// @scheme						http
// @BasePath					/
// @securityDefinitions.apiKey	Bearer
// @in							header
// @name						JWT
func main() {
	// Initialize logger immediately
	utils.InitGlobalLogger("info", "json")
	utils.LogInfo("Logger initialized for early application setup", nil)

	// Load configuration
	if err := loadConfig(); err != nil {
		utils.LogFatal("Failed to load configuration", logrus.Fields{
			"error": err.Error(),
		})
		return
	}

	// Initialize logger
	initLogger()

	// Initialize utilities
	utils.InitValidator()

	// Initialize the database
	if err := initDatabase(); err != nil {
		utils.LogFatal("Failed to initialize the database", logrus.Fields{
			"error": err.Error(),
		})
		return
	}

	// Start the server
	startServer()
}

// loadConfig loads application configuration
func loadConfig() error {
	utils.LogInfo("Loading application configuration", nil)
	if err := config.LoadEnv(); err != nil {
		utils.LogError(err, "Error loading configuration", nil)
		return err
	}
	utils.LogInfo("Configuration loaded successfully", nil)
	return nil
}

// initLogger initializes the global logger
func initLogger() {
	logLevel := config.AppConfig.Log.Level
	logFormat := config.AppConfig.Log.Format
	utils.InitGlobalLogger(logLevel, logFormat)
	utils.LogInfo("Logger initialized", logrus.Fields{
		"logLevel":  logLevel,
		"logFormat": logFormat,
	})
}

// initDatabase initializes the database connection
func initDatabase() error {
	utils.LogInfo("Initializing database", nil)
	if err := database.InitDatabase(); err != nil {
		utils.LogError(err, "Error initializing database", nil)
		return err
	}
	utils.LogInfo("Database initialized successfully", nil)
	return nil
}

// startServer configures and starts the HTTP server
func startServer() {
	serverPort := config.AppConfig.Server.Port
	utils.LogInfo("Starting server", logrus.Fields{
		"port": serverPort,
	})

	utils.RegisterMetrics()
    utils.ExposeMetrics()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverPort),
		Handler:      api.NewRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		utils.LogFatal("Error starting server", logrus.Fields{
			"error": err.Error(),
		})
	}
}
