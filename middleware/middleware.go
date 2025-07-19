package middleware

import (
	"net/http"
)

/*
	BlankMiddleware is an example of how to create a middleware function.

	We're taking a http.Handler and returning which is a next handler in the chain - (either the next middleware or the actual handler).

	We're returning the whole thing as a http.Handler so that it can be passed to mux.Handle().

	The function could be decorated with a dependency injection container if needed like so:
	func BlankMiddleware(c *container.Container) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}
*/

// BlankMiddleware is an example blank template for creating middleware
func BlankMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Stuff that should happen before the actual handler / or next middleware in the chain
		// Example: Checking headers, modifying request, etc.

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Stuff that should happen after the actual handler / or next middleware in the chain
		// Example: Logging, cleanup,etc.
	})
}

/*
 *
 *
 * Define middleware functions below or in separate files inside this package
 */
