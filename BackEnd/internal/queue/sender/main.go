package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ltovarm/DomoticX/BackEnd/internal/queue/handle"
	"github.com/streadway/amqp"
)

func main() {

	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()
	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()
	// With the instance and declare Queues that we can publish and subscribe to.
	// Queues must be declare from json.
	q, err := channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	log.Println(" > RabbitMQ connection sucessful.")

	// Create message in JSON format
	messages := make(chan handle.MessageJSON)
	// Get value from tcp payload
	go handle.GetValue(messages)

	for {
		message, more := <-messages // Reads message from channel (blocking until message received)
		if !more {
			break
		}
		queueData(message, q, channelRabbitMQ)
		log.Println(" > Id: ", message.Id, ". Msg published.")
	}
}

func queueData(message handle.MessageJSON, q amqp.Queue, channelRabbitMQ *amqp.Channel) {

	// Convert to JSON
	body, err := json.Marshal(message)
	if err != nil {
		log.Fatalf(" > Error converting to JSON: %v", err)
	}
	// Publish message to queue
	err = channelRabbitMQ.Publish(
		"",     // Empty exchange
		q.Name, // Queue name
		false,  // Do not wait for confirmation
		false,  // Do not require confirmation from clients
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf(" > Error publishing message: %v", err)
	}
}
