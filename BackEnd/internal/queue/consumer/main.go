package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	query "github.com/ltovarm/DomoticX/BackEnd/internal/query"
	"github.com/ltovarm/DomoticX/BackEnd/internal/queue/handle"
	"github.com/rs/cors"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

// // WebSocket Upgrader configuration
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any source (CORS)
	},
}

// Structure for handling WebSocket clients
type WebSocketClients struct {
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

// Global variable to store WebSocket clients
var wsClients = WebSocketClients{
	clients: make(map[*websocket.Conn]bool),
}

func main() {

	// Connect to RabbitMQ and set up WebSocket
	q, connectRabbitMQ, channelRabbitMQ, err := connectToRabbitMQ()
	if err != nil {
		panic(err)
	}
	defer closeApp(connectRabbitMQ, channelRabbitMQ)

	go connectWebSocket()

	// Consume messages from RabbitMQ
	msgs, err := channelRabbitMQ.Consume(
		q.Name, // queue name
		"",     // consumer label
		true,   // autoack
		false,  // exclusivity
		false,  // no waiting
		false,  // no local
		nil,    // arguments
	)
	if err != nil {
		log.Fatalf(" > Error when registering the consumer: %s", err)
	}

	log.Println(" > Dequeue message")

	// Waiting for messages
	for msg := range msgs {
		// Decode the JSON of the message
		var dataSlice handle.MessageJSON
		if err := json.Unmarshal(msg.Body, &dataSlice); err != nil {
			log.Printf(" > Error decoding JSON: %s", err)
		}

		// Processing the message
		insertJsonToTable(dataSlice, "temperatures")
		sendtoReact(msg.Body)
	}
}

// insertJsonToTable inserts JSON data into the specified SQL table
func insertJsonToTable(dataSlice handle.MessageJSON, sqltable string) {

	// Loop through the data slice
	for _, data := range dataSlice.Data {
		// Set-up connection
		my_db := query.NewDb()
		if err := my_db.ConnectToDatabaseFromEnvVar(); err != nil {
			log.Fatalf(" > Error connecting to db: %s\n", err)
		}

		// Convert to JSON
		body, err := json.Marshal(data)
		if err != nil {
			log.Println(" > Error converting JSON:", err)
		}

		// Insert into table
		if err := my_db.SendDataAsJSON(body, sqltable); err != nil {
			log.Fatalf(" > Error inserting data to db: %s\n", err)
		}
		my_db.CloseDatabase()
	}
}

// sendtoReact sends the message to all connected WebSocket clients
func sendtoReact(body []byte) {
	// Send data to all connected clients
	for client := range wsClients.clients {
		err := client.WriteMessage(websocket.TextMessage, body)
		if err != nil {
			log.Printf(" > Error sending message to a client: %s", err)
			// Optionally, you can remove the client from the list and close the connection.
			// If you don't want to close the connection, you can continue with the next client.
			wsClients.mutex.Lock()
			delete(wsClients.clients, client)
			wsClients.mutex.Unlock()
			client.Close()
		}
	}
}

// // connectToRabbitMQ connects to RabbitMQ and declares the queue
func connectToRabbitMQ() (q amqp.Queue, connectRabbitMQ *amqp.Connection, channelRabbitMQ *amqp.Channel, err error) {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err = amqp.Dial(amqpServerURL)
	if err != nil {
		log.Fatalln(err.Error())
		panic(err)
	}

	// Channel configuration
	channelRabbitMQ, err = connectRabbitMQ.Channel()
	if err != nil {
		log.Fatalf(" > Error opening a channel: %s", err)
	}

	// Queue configuration
	q, err = (*channelRabbitMQ).QueueDeclare(
		"QueueService1", // queue name
		true,            // durability
		false,           // self-deletion
		false,           // exclusivity
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf(" > Error in declaring the queue: %s", err.Error())
	}
	log.Println(" > Connection to RabbitMQ successful")

	return q, connectRabbitMQ, channelRabbitMQ, err
}

// // connectWebSocket sets up a WebSocket connection to send data to the frontend

func connectWebSocket() {

	// Establish WebSocket connection to send data to the frontend
	http.HandleFunc("/ws", handleWebSocket)

	// Configure CORS to allow requests from localhost:8080 (React)
	handler := cors.Default().Handler(http.DefaultServeMux)
	log.Println(" > Handler successful")
	log.Println(" > Server running at http://localhost:8080")
	http.ListenAndServe(":8080", handler)

}

// handleWebSocket handles WebSocket connections and reads incoming messages from the clients
func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(" > Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Add the client connection to the map
	wsClients.clients[conn] = true

	// Read incoming messages from the client (not used in this example)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Remove the client connection from the map when disconnected
			delete(wsClients.clients, conn)
			break
		}
	}
}

// closeApp closes the RabbitMQ connection and channel when the application is shutting down
func closeApp(conn *amqp.Connection, channel *amqp.Channel) {
	conn.Close()
	channel.Close()
	log.Println("App closed.")
}
