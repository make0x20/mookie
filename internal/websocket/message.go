package websocket

/*
   Message types and structure for the websocket package.
   Messages support both text and binary WebSocket frames while maintaining JSON structure.

   Default message types:
       MessageTypeError = "error" - Used for error responses

   Message structure:
       Type    - Application-level message type (e.g., "chat", "error")
       Payload - Message content as bytes
       Mode    - WebSocket frame type (text/binary)
       ClientID - Identifier of the sending client (set by server)

   Example usage:
       // Create and send a text message
       msg := Message{
           Type: "chat",
           Payload: []byte("Hello"),
           Mode: MessageModeText,
       }

       // Create and send a binary message
       msg := Message{
           Type: "data",
           Payload: []byte{1, 2, 3},
           Mode: MessageModeBinary,
       }
*/

// Message types
const (
	// Application message types
	MessageTypeError = "error"

	// Websocket message modes
	MessageModeText   = 1
	MessageModeBinary = 2
)

// Message structure
type Message struct {
	Mode     int    `json:"-"`
	Type     string `json:"type"`
	Payload  []byte `json:"payload"`
	ClientID string `json:"cid,omitempty"`
}
