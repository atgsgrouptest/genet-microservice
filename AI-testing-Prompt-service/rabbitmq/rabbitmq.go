package rabbitmq

import (
	"log"
	"os"
	"encoding/json"
	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
)

func InitRabbitMQ() {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	queue, err = channel.QueueDeclare(
		"request_ids", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

func PublishRequestID(requestID string , companyID string) error {
	data, _ := json.Marshal(map[string]string{
	"requestID":  requestID,
	"companyID":  companyID,
})

err := channel.Publish(
	"",           // exchange
	queue.Name,   // routing key
	false,
	false,
	amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	},
)
	return err
}

func Close() {
	channel.Close()
	conn.Close()
}
