package websocket

import (
	"errors"
	"sync"
)

/*
   How to use:
   1. Create a new Hub
   2. Add clients as they connect
   3. Use Broadcast() or SendToClients() to send messages
   4. Remove clients when they disconnect
   5. Close hub when shutting down

   Example:
       hub := websocket.NewHub()

       // Add new client
       client := websocket.NewClient("user123", conn, hub)
       hub.AddClient(client)

       // Broadcast to all clients
       hub.Broadcast(Message{
           Type: "announcement",
           Payload: []byte("Server starting"),
       })

       // Send to specific clients
       hub.SendToClients([]*Client{client1, client2}, Message{
           Type: "private",
           Payload: []byte("Hello"),
       })

       // Cleanup
       hub.Close()

   Notes:
   - Thread-safe client management
   - Supports broadcasting to all clients
   - Supports sending to specific clients
   - Handles client cleanup on disconnect
*/

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients []*Client
	mu      sync.RWMutex
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients: make([]*Client, 0),
	}
}

// AddClient adds a client to the hub.
func (h *Hub) AddClient(client *Client) error {
	if client == nil {
		return errors.New("client cannot be nil")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients = append(h.clients, client)
	return nil
}

// RemoveClient removes a client from the hub.
func (h *Hub) RemoveClient(client *Client) {
	if client == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, c := range h.clients {
		if c == client {
			h.clients = append(h.clients[:i], h.clients[i+1:]...)
			break
		}
	}
}

// Broadcast sends a message to all clients in the hub.
func (h *Hub) Broadcast(message Message) {
	h.mu.RLock()
	clients := make([]*Client, len(h.clients))
	copy(clients, h.clients) // Copy to avoid holding lock during send
	h.mu.RUnlock()

	for _, client := range clients {
		go func(c *Client) {
			c.Writer() <- message
		}(client)
	}
}

// SendToClients sends a message to a list of clients.
func (h *Hub) SendToClients(clients []*Client, message Message) {
    for _, client := range clients {
        go func(c *Client) {
            c.Writer() <- message
        }(client)
    }
}

// Close closes the hub and all clients.
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, client := range h.clients {
		client.Close()
	}
	h.clients = nil
}

// GetClients returns a list of clients in the hub.
func (h *Hub) GetClients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients
}
