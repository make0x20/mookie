package handlers

import (
	"mookie/internal/container"
	ws "mookie/internal/websocket"
	"mookie/templates/pages"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

/*
Define all the handler functions for the application here or in separate files inside the handlers package. Render the templates using templ's Render method.
*/
func Front() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		pages.Front().Render(r.Context(), w)
	}
}

func PostMessage(c *container.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get dependencies
		logger := c.MustGet("logger").(*slog.Logger)
		hub := c.MustGet("hub").(*ws.Hub)
		message := r.Header.Get("message")

		logger.Debug("received message", "message", message)

		// Check if the header exists and handle accordingly.
		if message == "" {
			http.Error(w, "message cannot be empty", http.StatusBadRequest)
			logger.Info("message cannot be empty")
			return
		}

		// Create a new websocket message
		wsMessage := ws.Message{
			Mode:    ws.MessageModeText,
			Type:    "message",
			Payload: []byte(message),
		}

		// Broadcast the message to all connected clients on the hub
		hub.Broadcast(wsMessage)

		// Respond with a 200 OK status
		w.WriteHeader(http.StatusOK)
	}
}

func BroadcastMessage(c *container.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get dependencies
		logger := c.MustGet("logger").(*slog.Logger)
		hub := c.MustGet("hub").(*ws.Hub)
		upgrader := c.MustGet("upgrader").(*websocket.Upgrader)

		// Upgrade the connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("failed to upgrade connection", "error", err)
			return
		}

		// Create a new client
		client := ws.NewClient("", conn, hub)

		// Add the client to the hub
		if err := hub.AddClient(client); err != nil {
			logger.Error("failed to add client", "error", err)
			conn.Close()
			return
		}

		// Start client
		if err := client.Start(); err != nil {
			logger.Error("failed to start client", "error", err)
			hub.RemoveClient(client)
			return
		}

		// Send connection message
		client.Writer() <- ws.Message{
			Type:    "connection",
			Payload: []byte("Connected to server"),
		}
	}
}
