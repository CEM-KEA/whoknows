package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/api"
	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
)

func main() {
	// Load application configuration from the .env file
	err := config.LoadEnv()
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		return
	}

	// Initialize the database
	err = database.InitDatabase()
	if err != nil {
		fmt.Printf("Error initializing database: %s\n", err)
		return
	}

	// Migrate the database schema to the latest version if the migrate flag is set to true
	if config.AppConfig.Database.Migrate {
		err = database.MigrateDatabase()
		if err != nil {
			fmt.Printf("Error migrating database: %s\n", err)
			return
		}
	}

	// Seed the database with initial data if the seed flag is set to true
	if config.AppConfig.Database.Seed {
		err = database.SeedData(database.DB, config.AppConfig.Database.SeedFilePath)
		if err != nil {
			fmt.Printf("Error seeding database: %s\n", err)
			return
		}
	}

	// initalize utils
	utils.InitValidator()

	// Create the router from the api package
	router := api.NewRouter()

	// Start the server
	serverPort := config.AppConfig.Server.Port
	fmt.Printf("Server is running on port: %d\n", serverPort)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err = server.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
