package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Telefonica/redis-vs-nats/measurement"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type rabbitMqMessagingService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewRabbitMqMessagingService() (MessagingService, error) {
	topic := "hello"

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err)
	}

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err)
	}

	return &rabbitMqMessagingService{
		conn: conn,
		ch:   ch,
		q:    q,
	}, nil
}

func (s *rabbitMqMessagingService) Publish(message *model.Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("fail to marshall message: %s", err)
	}

	return s.ch.Publish(
		"",       // exchange
		s.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        messageJSON,
		})
}

func (s *rabbitMqMessagingService) Subscribe(ready chan bool, size uint, ms *measurement.Measurements) error {
	msgs, err := s.ch.Consume(
		s.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to register a consumer")
	}

	ready <- true
	log.Println("listening")
	for d := range msgs {
		var message model.Message
		if err = json.Unmarshal(d.Body, &message); err != nil {
			log.Printf("error unmarshaling message: %s", err)
		} else {
			if message.ID == 1 {
				ms.Start()
				ms.SetFirstReceiving(*message.SentAt)
			}
			ms.AddMeasurement(time.Since(*message.SentAt))
			if message.ID == size {
				ms.RegisterLastReceiving()
				break
			}
		}
	}
	ms.Stop()
	return nil
}
