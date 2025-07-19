package routes

import (
	"mookie/handlers"
	"mookie/internal/container"
	"mookie/middleware"
	"net/http"
)

/*
Define all the routes for the application here. Chain middleware together
with the middleware.Chain function, and use the http.HandlerFunc function
to convert your handler functions to http.Handler types.
*/
func Setup(c *container.Container) http.Handler {
	// Setup middlewares
	// Default middleware chain - pass the dependency container
	defaultChain := middleware.DefaultChain(c)

	// Create a new ServeMux router
	mux := http.NewServeMux()

	// Define routes - replace with your own
	// Load frontpage
	mux.Handle("GET /", defaultChain(
		http.HandlerFunc(handlers.Front())),
	)

	// Post message
	mux.Handle("POST /post-message", defaultChain(
		http.HandlerFunc(handlers.PostMessage(c))),
	)

	// Websocket message stream
	mux.Handle("GET /ws/message-stream", defaultChain(
		http.HandlerFunc(handlers.BroadcastMessage(c))),
	)

	// Serve static files from static folder as /static/*
	fs := http.FileServer(http.Dir("static"))
	staticHandler := http.StripPrefix("/static/", fs)
	mux.Handle("GET /static/", defaultChain(staticHandler))

	return mux
}
