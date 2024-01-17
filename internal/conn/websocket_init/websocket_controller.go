package websocketinit

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang-websocket-amqp-chatapp/internal/conn/mq"
	"golang-websocket-amqp-chatapp/internal/types"
)

func SendToUserByID(sender *websocket.Conn, senderID string, targetID string, message []byte) {
	clientsMtx.Lock()
	defer clientsMtx.Unlock()

	if targetID == senderID {
		// Sender and target are the same user
		errMsg := "You cannot send a message to yourself."
		sendMessage(sender, []byte(errMsg))
		return
	}

	for clientConn, clientID := range clients {
		if clientID == targetID {
			// Target user found, send the message
			sendMessage(clientConn, message)
			deliveredMsg := "Delivered.✔️✔️"
			sendMessage(sender, []byte(deliveredMsg))
			return
		}
	} // Target user not found
	sentMsg := "Delivered.✔️"
	sendMessage(sender, []byte(sentMsg))
	unavailableUser(senderID, targetID, message)

}

func SendLateToUserByID(sender *websocket.Conn, senderID string, targetID string, message []byte) {
	clientsMtx.Lock()
	defer clientsMtx.Unlock()

	for clientConn, clientID := range clients {
		if clientID == targetID {
			// Target user found, send the message
			sendMessage(clientConn, message)

			return
		}
	}

}

// sendToAll sends a message to all connected clients except the sender.
func sendToAll(message types.Message) {
	// Lock the mutex to ensure safe concurrent access to the clients map
	clientsMtx.Lock()
	defer clientsMtx.Unlock()

	// Iterate through all connected clients
	for client := range clients {

		// Check if the current client is not the sender
		if client != message.Sender {
			// Send the message to the current client
			sendMessage(client, message.Message)
			sentMsg := "Delivered.✔️"
			sendMessage(message.Sender, []byte(sentMsg))
		}
	}
}

// unavailableUser handles the scenario when a user is unavailable and a message needs to be sent.
func unavailableUser(sender string, target string, message []byte) {
	// Convert the message bytes to a string
	strMsg := string(message)

	// Create a map representing the message data
	data := types.MessageData{
		TargetID: target,
		Message:  strMsg,
		SenderID: sender,
	}

	// Convert the map to a JSON string
	dataString, err := json.Marshal(data)
	if err != nil {
		// Print an error message if there's an issue with JSON marshalling
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Convert the JSON byte slice to a string
	jsonString := string(dataString)

	// Publish the message to the "chat" queue using the MQ (Message Queue) service
	mq.PublishMSG("chat", jsonString)
}
