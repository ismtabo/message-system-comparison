package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/streadway/amqp"
)

var filename string = "../../json/100k.json"
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
	defer conn.Close()
	defer ch.Close()

	start := time.Now()

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("fail open json file: %s", err)
	}

	decoder := json.NewDecoder(jsonFile)

	// Read opening file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail to decode token: %s", err)
	}

	var message model.Message
	for decoder.More() {
		err := decoder.Decode(&message)
		if err != nil {
			log.Fatalf("fail to decode message: %s", err)
		}

		now := time.Now()
		message.SentAt = &now

		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("fail to marshall message: %s", err)
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        messageJSON,
			})
		if err != nil {
			log.Fatalf("fail to marshall message: %s", err)
		}
	}

	// Close the file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail decoding token: %s", err)
	}

	elapsed := time.Since(start)
	log.Printf("Nats Sender took %s\n", elapsed)
}
