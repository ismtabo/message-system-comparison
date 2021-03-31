package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	ch      *amqp.Channel
	q       amqp.Queue
	size    int64 = 100_000
	err     error
	n       int64
	Ex, Ex2 float64
)

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

	go func(last chan bool, expectedMessages int64) {
		defer close(last)
		for d := range msgs {
			var message model.Message
			if err = json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("error unmarshaling message: %s", err)
			} else {
				x := time.Since(*message.SentAt).Seconds()
				n++
				Ex += x
				Ex2 += x * x
				if n >= expectedMessages {
					break
				}
			}
		}
	}(forever, size)

	log.Println("listening")
	<-forever

	fN := float64(n)
	fEx := float64(Ex)
	fEx2 := float64(Ex2)

	mean := fEx / fN
	variance := (fEx2 - (fEx*fEx)/fN) / (fN - 1)
	log.Printf("Messages received: %d Latency: mean %gs, variance %gs Throughput: %g msg/s", n, mean, variance, 1/mean)
}
