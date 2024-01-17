package mq

import (
	"golang-websocket-amqp-chatapp/internal/types"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"os"
)

var conn *amqp.Connection

// Global channel
var GlobalChannel = make(chan string)

func InitAMPQ() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("AMPQURI")

	var err error
	conn, err = amqp.Dial(uri)
	if err != nil {
		log.Println(uri)
		panic(err)
	}
	log.Println("Successfully connected to AMQP")
	createQueues("chat")
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// createQueues checks if the specified queue already exists. If it does, it logs a message
// and does nothing. If the queue doesn't exist, it declares the queue and logs a message
// indicating that the queue has been created.
func createQueues(qName string) {
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create an amqpChannel")
	defer amqpChannel.Close()

	// Check if the queue already exists
	_, err = amqpChannel.QueueDeclarePassive(qName, true, false, false, false, nil)
	if err == nil {
		// Queue already exists, do nothing
		log.Printf("Queue %s is already here!", qName)
		return
	}

	// If the queue does not exist, declare it
	_, err = amqpChannel.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		log.Printf("Uh-oh! Error declaring queue %s: %v", qName, err)
	} else {
		log.Printf("Queue %s has been created successfully!", qName)
	}
}

func PublishMSG(qName string, data interface{}) {
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	err = amqpChannel.Publish("", qName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(data.(string)),
	})
	handleError(err, "Error publishing a message to the queue")
}

func Consumer(id string) (bool, types.MessageData) {
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create an amqpChannel")
	defer amqpChannel.Close()
	// Use a channel to communicate the result from the goroutine
	// Use a channel to communicate the result and decoded data from the goroutine
	resultChan := make(chan struct {
		success bool
		data    types.MessageData
	})
	queueNames := []string{"chat", "request"} // Add the queue names you want to check
	// Initialize the newBalance channel
	for _, queueName := range queueNames {
		queueName := queueName // Create a local variable and assign the value of queueName
		queue, err := amqpChannel.QueueDeclare(queueName, true, false, false, false, nil)
		handleError(err, "Could not declare queue")

		messages, err := amqpChannel.Consume(queue.Name, "", false, false, false, false, nil)
		handleError(err, "Could not register consumer")

		go func(qName string) {
			for d := range messages {
				log.Printf("Received a message from queue %s", d.Body)

				if qName == "chat" {

					// Assuming 'jsonString' is a JSON-encoded string
					var decodedData types.MessageData
					err := json.Unmarshal(d.Body, &decodedData)
					if err != nil {
						fmt.Println("Error unmarshalling JSON:", err)
						continue
					}

					fmt.Println("Decoded Map:", decodedData)
					fmt.Println("my id:", id)
					// Perform your condition check here
					if decodedData.TargetID == id {
						fmt.Println("This is what i sent:", decodedData.Message)

						// Acknowledge the message only if the condition is met
						err := d.Ack(true)
						if err != nil {
							log.Printf("Error acknowledging message: %v", err)
						}
						resultChan <- struct {
							success bool
							data    types.MessageData
						}{true, decodedData}
						return
					} else {
						// Handle the case when the condition is not met (e.g., requeue or reject)
						err := d.Nack(true, true)
						if err != nil {
							log.Printf("Error rejecting message: %v", err)
						}
						resultChan <- struct {
							success bool
							data    types.MessageData
						}{false, types.MessageData{}}
						return
					}
				}
			}
		}(queueName)
	}

	// Wait for the result from the goroutine
	result := <-resultChan

	// Return the final result
	return result.success, result.data

}
