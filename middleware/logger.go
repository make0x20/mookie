package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

// LoggerMiddleware logs the request information
// It should be the first middleware in the chain
func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get current time for request duration
			start := time.Now()

			// Generate and set request ID
			requestID := uuid.New().String()
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			r = r.WithContext(ctx)
			w.Header().Set("X-Request-ID", requestID)

			// Get real IP if behind proxy
			realIP := r.Header.Get("X-Real-IP")
			if realIP == "" {
				realIP = r.Header.Get("X-Forwarded-For")
			}
			if realIP == "" {
				realIP = r.RemoteAddr
			}

			// Call the next middleware or final handler in the chain
			next.ServeHTTP(w, r)

			// Get query parameters
			var queryParams string
			if r.URL.RawQuery != "" {
				queryParams = "?" + r.URL.RawQuery
			}

			logger.Info("http request",
				"request_id", requestID,
				"method", r.Method,
				"protocol", r.Proto,
				"duration", time.Since(start).String(),
				"ip", realIP,
				"host", r.Host,
				"path", r.URL.Path+queryParams,
				"user_agent", r.UserAgent(),
				"referer", r.Referer(),
			)
		})
	}
}
