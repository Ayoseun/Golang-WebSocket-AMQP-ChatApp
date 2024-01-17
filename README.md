# Golang-WebSocket-AMQP-ChatApp

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Golang Version](https://img.shields.io/badge/golang-v1.17-blue)
![License](https://img.shields.io/badge/license-MIT-blue)

**Author:** Ayoseun

## Overview
A simple go server, Engineered for real-time bidirectional communication using WebSocket and persistent messaging with AMQP. Efficient concurrency model and heartbeat mechanisms ensure responsiveness.


## Description 
The server is structured around the WebSocket protocol, establishing a full-duplex communication channel. This bidirectional nature allows for simultaneous data flow between the server and clients, optimizing responsiveness and real-time interaction.

For persistent communication, I integrated the Advanced Message Queuing Protocol (AMQP). This decision ensures message durability and reliable delivery, even when clients experience temporary disconnections. Clients reconnecting can retrieve any missed messages, creating a resilient communication system.

Concurrently managing read and write operations is crucial for efficiency. Hence, I implemented a well-thought-out concurrency model utilizing goroutines and channels. This approach enables parallel execution of communication tasks, avoiding bottlenecks and enhancing scalability.

To maintain the health of WebSocket connections, I incorporated heartbeat mechanisms using Ping and Pong. This proactive strategy prevents connection staleness by regularly checking and managing timeouts.

In summary, the server architecture combines WebSocket for real-time bidirectional communication, AMQP for persistent messaging, and an efficient concurrency model. This engineering design ensures a responsive, reliable, and scalable system for facilitating dynamic interactions within a distributed environment.


This Golang server powers a robust real estate communication infrastructure, fostering seamless interactions between administrators, landlords, and tenants. Leveraging WebSockets and MQTT protocols, it ensures efficient, real-time communication, making it an ideal foundation for comprehensive real estate solutions.

## Features

- **Real-Time Communication:** Utilizes WebSockets and MQTT to facilitate instantaneous messaging.
- **User Roles:** Supports distinct user roles—admin, landlord, and tenant—for tailored communication.
- **Scalability:** Engineered for scalability, ensuring optimal performance as user interactions grow.
- **Concurrent Handling:** Leverages Golang's concurrency model for managing multiple simultaneous requests.
- **Customizable:** Easily adaptable for various real estate use cases, from managing property listings to tenant inquiries.

## Additional Features

- **User Connection Management:**
Your code doesn't currently handle user connections and disconnections explicitly. You need a mechanism to manage and identify users based on their connections.

- **Message Delivery:**
When a user connects (Consumer starts listening), you may want to check for any pending messages for that user in the queue and deliver them.

-**Handling Disconnects:**
Implement a way to gracefully handle user disconnections, and possibly remove their connection information from your internal data structures.

-**Acknowledgment Mechanism:**
Implement an acknowledgment mechanism so that a sender knows when a message has been successfully delivered to the recipient.


# Get started

create .env file in the root directory and add AMQP url and AUTHHEADER url

- To run this in local environment during development
  
```
go run cmd/main.go

```

- For production it is better to use the make file 

```
make
```

- connect to websocket using /ws, you will need to set an id because the id is used to map each client to the client hashmap

```shell
ws://localhost:8080/ws?id=1113
```

-payload format

```json
{
    "target_id":"48trwSQ00",
    "message":"yo, admin chat me up",
    "sender_id":"18trwSQ00"
}

```

1. target_id: this is the id of the recipient
2. message : This is the messsage that is been sent
3. sender_id : This should be the same id as the id in the parameter