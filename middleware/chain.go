package middleware

import (
	"mookie/internal/container"
	"log/slog"
	"net/http"
)
/*
	I wrote this comment cause I realize it's a bit confusing to get the hang of this at first :)

	Chain does the magic of chaining middlewares together
	Basically, it takes a http.Handler and a list of middlewares functions

	It loops through each middleware and runs the middleware function which returns a http.Handler
	that is wrapped with the current middleware function, then it moves to the next middleware function
	and does the same thing until the end.

	It basically becomes a handler that has all the middlewares applied in order.

	It executes like so: middleware3(middleware2(middleware1(handler)))
*/

// Chain applies middlewares in order
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

/*
	DefaultChain essentially returns a chain of middlewares. It's a function that takes http.Handler and returns http.Handler.

	This is a function that we use to pass the handler.

	We pass middleware functions inside the Chain which returns a http.Handler (the whole chain) which has the correct http.HandlerFunc signature and
	can be passed to mux.Handle().

	The outer function is a decorator so that we can pass the dependency injection container.
*/

// DefaultChain is a default chain of middlewares
func DefaultChain(c *container.Container) func(http.Handler) http.Handler {
	logger := c.MustGet("logger").(*slog.Logger)
	return func(h http.Handler) http.Handler {
		return Chain(h,
			LoggerMiddleware(logger),
			// BlankMiddleware,
		)
	}
}

/*
 *
 *
 * A place for more middleware chains
 */
