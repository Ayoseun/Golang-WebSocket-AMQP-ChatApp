package main

import (
	"log"
	"github.com/joho/godotenv"
	"golang-websocket-amqp-chatapp/internal/conn"
)

func main() {

	// load .env file
	goDotEnvVariable()
	conn.Setup()

}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable() {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}
