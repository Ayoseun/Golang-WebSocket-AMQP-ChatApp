package types

import "github.com/gorilla/websocket"

// Message struct to hold information about the sender and the actual message
type Message struct {
	Sender  *websocket.Conn // WebSocket connection of the sender
	Message []byte           // The content of the message
}

// MessageData struct to represent the structure of a message with specific fields
type MessageData struct {
	SenderID string `json:"sender_id"` // ID of the sender
	TargetID string `json:"target_id"` // ID of the target recipient
	Message  string `json:"message"`    // Content of the message
}

// User represents a connected user with an authentication token and optional admin status.
type User struct {
	Conn   *websocket.Conn // WebSocket connection of the user
	Token  string           // Authentication token for the user
	IsAdmin bool            // Flag indicating whether the user has admin privileges
}
