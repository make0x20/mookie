package websocket

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
)

/*
   Package websocket provides a WebSocket client implementation with support for
   text/binary frames and JSON message handling.

   How to use:
   1. Create a new Client with an ID, WebSocket connection, and Hub
   2. Start the client (starts read/write pumps)
   3. Use Reader() and Writer() channels to communicate
   4. Close when done

   Example basic usage:
       hub := websocket.NewHub()

       // When websocket connects
       client := websocket.NewClient("user123", conn, hub)
       if err := hub.AddClient(client); err != nil {
           log.Println("Error adding client:", err)
           return
       }
       if err := client.Start(); err != nil {
           log.Println("Error starting client:", err)
           hub.RemoveClient(client)
           return
       }

   Example sending messages:
       // Send text message
       client.Writer() <- Message{
           Type:    "chat",
           Payload: []byte("Hello"),
           Mode:    MessageModeText,
       }

       // Send binary message
       client.Writer() <- Message{
           Type:    "data",
           Payload: []byte{1, 2, 3},
           Mode:    MessageModeBinary,
       }

   Example receiving messages:
       for msg := range client.Reader() {
           switch msg.Type {
           case "chat":
               // Handle chat message
           case "data":
               // Handle data message
           case MessageTypeError:
               // Handle error message
           }
       }

   Notes:
   - Supports both text and binary WebSocket frames
   - Automatically handles WebSocket control frames (ping/pong)
   - Thread-safe message handling through channels
   - Automatic cleanup on connection close
   - JSON message encoding/decoding
   - Integrates with Hub for broadcast capabilities
   - Buffered channels (256 messages)
*/

// Client represents a WebSocket client
type Client struct {
	ID      string
	conn    *websocket.Conn
	send    chan Message
	receive chan Message
	hub     *Hub
}

// NewClient creates a new WebSocket client
func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:      id,
		conn:    conn,
		send:    make(chan Message, 256),
		receive: make(chan Message, 256),
		hub:     hub,
	}
}

// Start the client read/write pumps
func (c *Client) Start() error {
	if c.conn == nil {
		return errors.New("connection not initialized")
	}
	go c.readPump()
	go c.writePump()
	return nil
}

// Close the client connection
func (c *Client) Close() {
	if c.conn == nil {
		return
	}

	c.conn.Close()
	close(c.send)
	close(c.receive)
}

// Reader returns the receive channel
func (c *Client) Reader() <-chan Message {
	return c.receive
}

// Writer returns the send channel
func (c *Client) Writer() chan<- Message {
	return c.send
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.RemoveClient(c)
		c.Close()
	}()

	for {
		messageType, payload, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		if err := c.handleMessage(messageType, payload); err != nil {
			return
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(messageType int, payload []byte) error {
	switch messageType {
	case websocket.TextMessage, websocket.BinaryMessage:
		return c.handleDataMessage(messageType, payload)

	case websocket.CloseMessage:
		return errors.New("close message received")

	case websocket.PingMessage:
		return c.conn.WriteMessage(websocket.PongMessage, nil)

	case websocket.PongMessage:
		return nil
	}

	return nil
}

// handleDataMessage processes incoming data messages
func (c *Client) handleDataMessage(messageType int, payload []byte) error {
	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		c.send <- Message{
			Type:    MessageTypeError,
			Payload: []byte("Invalid message"),
			Mode:    messageType,
		}
		return nil
	}

	msg.ClientID = c.ID
	msg.Mode = messageType
	c.receive <- msg
	return nil
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	for msg := range c.send {
		switch msg.Mode {
		case MessageModeBinary:
			data, err := json.Marshal(msg)
			if err != nil {
				return
			}
			err = c.conn.WriteMessage(websocket.BinaryMessage, data)
		default:
			// Default to text message mode
			data, err := json.Marshal(msg)
			if err != nil {
				return
			}
			err = c.conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
