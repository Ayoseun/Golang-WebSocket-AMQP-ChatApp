package websocketinit

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// WebSocket Upgrader configuration
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allowing any origin for WebSocket connections
			return true
		},
	}
	 lastMsgMutex sync.Mutex

	// Map to keep track of connected clients, associating each WebSocket connection with a user ID
	clients    = make(map[*websocket.Conn]string)
	clientsMtx sync.Mutex // A lock to prevent simultaneous access to the clients map
)

