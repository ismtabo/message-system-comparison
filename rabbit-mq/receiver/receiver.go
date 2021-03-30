package main

import (
	"encoding/json"
	"log"

	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel
var q amqp.Queue
var err error

func init() {
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err)
	}

	q, err = ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err)
	}
}

func main() {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var message model.Message
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Error unmarshaling message: %s", err)
			}
			log.Printf("Received a message: %v\n", message)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
