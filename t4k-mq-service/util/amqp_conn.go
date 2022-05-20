package util

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func NewAmqpConnection() *amqp.Connection {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		log.Fatalf("failed to connect to MQ: %v", err)
	}
	return conn
}
