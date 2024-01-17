package mqutil

import (

	"fmt"
	

	"golang-websocket-amqp-chatapp/internal/conn/mq"
	"golang-websocket-amqp-chatapp/internal/types"
)

func AwaitingMsgs(id string) types.MessageData {
    // Use a channel to signal the completion of message processing
    done := make(chan struct{})
    defer close(done)

    var msgData types.MessageData
	awaitingMsg, data := mq.Consumer(id)
        if awaitingMsg {
            msgData = data
            fmt.Println("Decoded Map from util:", msgData,awaitingMsg)
        }

    // Return the decoded message data
    return msgData
}

