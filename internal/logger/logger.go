package logger

import (
	"io"
	"log/slog"
	"os"
)

/*
   Package logger provides a simple structured logging setup using slog.
   It supports multiple writers and configurable log levels.

   How to use:
   1. Create a new logger with desired log level
   2. Optionally provide additional writers (e.g., file, network)
     - This allows logging to multiple destinations whether it's stdout, file, network, or other
   3. Use standard slog methods for logging

   Example with stdout only:
       logger := logger.New(slog.LevelInfo)
       logger.Info("Server starting", "port", 8080)

   Example with file and stdout:
       file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
       if err != nil {
           log.Fatal(err)
       }
       logger := logger.New(slog.LevelDebug, file)

       // Logs to both stdout and file
       logger.Debug("Config loaded", "config", cfg)
       logger.Error("Connection failed", "error", err)

   Notes:
   - Always writes to stdout
   - Additional writers are optional
   - Nil writers are filtered out
   - Uses slog's text handler for readable output
*/

// New creates a new logger with the given log level and io.writer
func New(level slog.Level, writers ...io.Writer) *slog.Logger {
	// Always include stdout writer
	validWriters := []io.Writer{os.Stdout}

	// Filter out nil writers
	for _, w := range writers {
		if w != nil {
			validWriters = append(validWriters, w)
		}
	}

	// Combine writers into multiwriter
	mWriter := io.MultiWriter(validWriters...)
	// Set log level
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Create new logger
	return slog.New(slog.NewJSONHandler(mWriter, opts))
}
