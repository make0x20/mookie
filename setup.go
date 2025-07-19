// setup.go
package main

import (
	"context"
	"fmt"
	ws "github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"mookie/config"
	"mookie/internal/container"
	"mookie/internal/db"
	"mookie/internal/db/sqlc"
	"mookie/internal/logger"
	"mookie/internal/websocket"
	"net/http"
	"os"
)

// setupDependencies initializes and registers all application dependencies.
// Add or modify dependencies here as needed for your project.
func setupDependencies(configPath *string) (*container.Container, error) {
	// Create a new dependency injection container
	container := container.New()

	// Load the config
	cfg := setupConfig(configPath)
	container.Register("config", cfg)

	// Setup logger
	logger := setupLogger(cfg)
	container.Register("logger", logger)

	// Debug log config
	logger.Debug("Loaded config", "config", cfg)

	// Set up database
	db, err := db.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	container.Register("db", db)

	// Set up websocket hub
	hub := websocket.NewHub()
	container.Register("hub", hub)
	// Set up websocket upgrader - allow all origins for now
	upgrader := &ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	container.Register("upgrader", upgrader)

	return container, nil
}

// setupLogger is a helper function that creates a new logger with the specified configuration - log file and log level
func setupLogger(cfg *config.Config) *slog.Logger {
	var file *os.File
	err := error(nil)

	// If a log file is specified, open it, otherwise log to stdout only
	if cfg.LogFile != "" {
		file, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}
	}

	logLevel := slog.LevelInfo
	// If the log level is debug, set it to debug otherwise leave it as info
	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}

	return logger.New(logLevel, file)
}

// setupConfig is a helper function that loads the configuration from the specified path
func setupConfig(path *string) *config.Config {
	cfg, err := config.NewWithPath(*path)
	if cfg == nil {
		log.Fatalf("error loading config: %v", err)
	}

	return cfg
}

// initDB initialized the db with predefined content - e.g. creating an admin user
func initDB(c *container.Container) {
	cfg := c.MustGet("config").(*config.Config)
	dbPath := cfg.DatabasePath

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	queries := sqlc.New(database)
	ctx := context.Background()

	// Check if admin user already exists
	_, err = queries.GetUserByUsername(ctx, "admin")
	if err == nil {
		fmt.Println("Admin user already exists, skipping creation")
		return
	}

	// Admin user doesn't exist, create it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created admin user: %+v\n", user)
}
