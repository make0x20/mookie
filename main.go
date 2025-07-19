package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"mookie/config"
	"mookie/routes"
	"net/http"
)

/*
Application structure:
	- main.go: Entry point of the application
	- setup.go: Define dependencies and set up the application
	- config/: Define configuration
	- handlers/: Define route handlers
	- internal/: Internal packages - should not be modified
		- container/: Simple dependency injection container system
		- cron/: Simple package to register cron jobs and run at specified intervals
		- db/: Database setup and connection - SQLite + sqlc
		- logger/: Structured logging setup using slog, allows multiple writers
		- websocket/: Simple websocket abstraction layer using Gorilla Websocket as the underlying library
	- middleware/: Define middleware
	- routes/: Define routes
	- static/: Static files
	- templates/: HTML templates using TEMPL template engine
	- services/: Suggested location for custom business logic

Application flow:
	1. Parse command line flags
	2. Set up dependencies
		- Load config
		- Set up logger
		- Set up database
		- Set up websocket hub and upgrader
	3. Set up routes and pass the container to the routes setup function
		- Routes define route handlers and middleware
	4. Start the server
*/

func main() {
	// Parse command line flags - define your own flags here if needed
	configPath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	// Set up dependencies - inside setup.go
	container, err := setupDependencies(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Get logger and config from the dependency container
	cfg := container.MustGet("config").(*config.Config)
	logger := container.MustGet("logger").(*slog.Logger)

	// Initialize database
	initDB(container)

	// Setup routes and pass the dependency container
	r := routes.Setup(container)

	addr := fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.Port)
	// Start the web server
	logger.Info("Starting server", "address", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}


