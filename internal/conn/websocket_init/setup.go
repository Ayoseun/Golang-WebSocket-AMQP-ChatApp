package websocketinit

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"golang-websocket-amqp-chatapp/internal/conn/mq/mq_util"
	"golang-websocket-amqp-chatapp/internal/types"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	defer func() {
		// Cleanup: Remove the client from the map and close the connection
		clientsMtx.Lock()
		delete(clients, conn)
		clientsMtx.Unlock()

		conn.Close()
	}()

	// Extract the 'id' parameter from the URL
	id := r.URL.Query().Get("id")
	if id == "" {
		// If 'id' is not provided, send a forbidden message and disconnect
		sendMessage(conn, []byte("Forbidden: Please provide an 'id' parameter"))
		return
	}

	clientsMtx.Lock()

	// Check if the client is already in the map
	_, exists := clients[conn]
	if exists {
		clientsMtx.Unlock()
		// If the client is already in the map, do nothing
		return
	}

	// Add the client to the map and send a connected message
	clients[conn] = id
	clientsMtx.Unlock()

	// Send a connected message
	sendMessage(conn, []byte("Connected"))
	go lastMsg(conn, id)

	// Handle incoming messages from the client
	for {
		_, clientMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client:", err)
			break
		}

		// Unmarshal the client's message
		var msg types.MessageData
		if err := json.Unmarshal(clientMsg, &msg); err != nil {
			sendMessage(conn, []byte("Bad payload"))
			continue
		}

		// Handle the client's message
		handleClientMessage(conn, msg)
	}
}


// handleClientMessage processes a message received from a client.
func handleClientMessage(sender *websocket.Conn, msg types.MessageData) {
	// Switch based on the sender's ID
	switch msg.SenderID {
	case "2024Ayoseun":
		// If the sender is a specific user, perform additional validation
		if msg.TargetID == "ayosinclusiveness" {
			// If the additional validation passes, broadcast the message to all connected clients
			sendToAll(types.Message{Sender: sender, Message: []byte(msg.Message)})
		} else {
			SendToUserByID(sender, msg.SenderID, msg.TargetID, []byte(msg.Message))
		}
	default:
		// For other senders, send the message to the specified target user
		SendToUserByID(sender, msg.SenderID, msg.TargetID, []byte(msg.Message))
	}
}


// sendMessage sends a message to a WebSocket connection.
func sendMessage(conn *websocket.Conn, message []byte) {
	// Attempt to write the message to the WebSocket connection
	err := conn.WriteMessage(websocket.TextMessage, message)

	// Check if there was an error writing to the connection
	if err != nil {
		log.Println("Error writing to client:", err)

		// Safely remove the client from the map and close the connection
		clientsMtx.Lock()
		delete(clients, conn)
		clientsMtx.Unlock()

		conn.Close()
	}
}


// lastMsg retrieves the last message for a specific user and sends it to the provided WebSocket connection.
func lastMsg(sender *websocket.Conn, id string) {
	lastMsgMutex.Lock()
	defer lastMsgMutex.Unlock()
	// Retrieve the last message for the specified user from the message queue
	data := mqutil.AwaitingMsgs(id)

	// Check if the retrieved message is intended for the specified user
	if id == data.TargetID {
		// Send the late message to the user using the SendLateToUserByID function
		SendLateToUserByID(sender, data.SenderID, data.TargetID, []byte(data.Message))
	}
}

// Rains initializes the WebSocket server and starts listening for connections.
func RouteSetup() {
	// Register the WebSocket endpoint ("/ws") with the authMiddleware to handle authentication
	http.HandleFunc("/ws", authMiddleware(handleWebSocket))

	// Log the server start-up message
	log.Println("WebSocket server listening on :8080")

	// Start the HTTP server and handle incoming WebSocket connections
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		// Log and exit the application if there's an error starting the server
		log.Fatal(err)
	}
}
