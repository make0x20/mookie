# Mookie Go starter project

A personal minimal starting point to build Go web applications and microservices.

## Features

- HTML templating with [TEMPL](https://templ.guide/)
- Middleware chain system
- Structured logging with slog
- Configuration via TOML and environment variables
- sqlc for database querying
- WebSocket support
- Cron job scheduling
- Dependency injection container
- Static file serving

## Structure

- main.go: Entry point of the application
- setup.go: Define dependencies and set up the application
- config/: Define configuration
- handlers/: Define route handlers
- internal/: Internal packages - should not be modified
	- container/: Simple dependency injection container system
	- cron/: Simple package to register cron jobs and run at specified intervals
	- logger/: Structured logging setup using slog, allows multiple writers
	- websocket/: Simple websocket abstraction layer using Gorilla Websocket as the underlying library
    - db/: Simple sqlite wrapper - combined with sqlc
- middleware/: Define middleware
- routes/: Define routes
- static/: Static files
- templates/: HTML templates using TEMPL template engine
- services/: Suggested location for custom business logic

## Quick start

### Prerequisites

- Install TEMPL: `go install github.com/a-h/templ/cmd/templ@latest`
- Install sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### Create a new project

- Clone this repository and cd into it
- Run `./rename-project.sh mookie <new-project-name>`
    - This will just rename mookie to your new project name and correct the imports
- Create `config.toml` via `cp config.toml.example config.toml`
- Run `go mod tidy` to install dependencies
- Run `templ generate` to generate initial compiled templates
- Run `go run .` to start the server
- Re-run `sqlc generate` whenever you change SQL queries, to regenerate the sqlc code

Optional:

- Install [Air](https://github.com/air-verse/air) and run `air` which will recompile and restart the server on file changes
