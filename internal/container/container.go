package container

import (
	"fmt"
	"sync"
)

/*
   Package container provides a simple dependency injection container for managing
   application-wide services and dependencies.

   How to use:
   1. Create a new Container
   2. Register services with unique names
   3. Retrieve services using Get (with error handling) or MustGet (panics on error)
   4. Type assert retrieved services to their concrete types

   Example basic usage:
       // Create container
       container := container.New()

       // Register core services
       container.Register("logger", slog.New(...))
       container.Register("db", db.Open(...))
       container.Register("config", config.Load())

       // Get service with error handling
       db, err := container.Get("db")
       if err != nil {
           log.Fatal(err)
       }
       dbInstance := db.(*db.DB)

       // Get service with panic on error and assert type (*slog.Logger in this case)
       logger := container.MustGet("logger").(*slog.Logger)

   Example in web application:
       func main() {
           container := container.New()

           // Register all dependencies
           container.Register("config", cfg)
           container.Register("logger", logger)
           container.Register("db", db)

           // Pass the container with dependencies to the router setup
           r := routes.Setup(container)
           http.ListenAndServe(":8080", r)
       }

   Notes:
   - Thread-safe
   - Services are stored as interface{} (any) which supports any dependency type
   - Type assertion required when retrieving services
   - Register will overwrite existing services with same name
   - MustGet panics if service not found
*/

// Container is a dependency injection container
type Container struct {
	services map[string]any
	mu       sync.RWMutex
}

// New creates a new dependency container
func New() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

// Register a service by name
func (c *Container) Register(name string, service any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get a service by name
func (c *Container) Get(name string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service, nil
}

// Type-safe getters
func (c *Container) MustGet(name string) any {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}
