package conn

import (
	"golang-websocket-amqp-chatapp/internal/conn/mq"
	"golang-websocket-amqp-chatapp/internal/conn/websocket_init"
)



// Setup initializes the necessary connections and services.
func Setup() {
	// Initialize the AMPQ (Advanced Message Queuing Protocol) connection
	mq.InitAMPQ()

	// Start the WebSocket server using the WebSocketInit package
	websocketinit.RouteSetup()
}
