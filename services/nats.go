package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Telefonica/redis-vs-nats/measurement"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type natsMessagingService struct {
	client *nats.Conn
	topic  string
}

func NewNatsMessagingService() (MessagingService, error) {
	client, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, errors.Wrap(err, "fail connecting to nats")
	}
	return &natsMessagingService{client: client, topic: "message"}, nil
}

func (s *natsMessagingService) Publish(message *model.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("fail marshaling message: %s", err)
	}
	return s.client.Publish(s.topic, data)
}

func (s *natsMessagingService) Subscribe(ready chan bool, size uint, ms *measurement.Measurements) error {
	last := make(chan bool)

	cb := func(last chan bool, expectedMessages uint) func(*nats.Msg) {
		return func(m *nats.Msg) {
			if ms.N() < expectedMessages {
				var message model.Message
				if err := json.Unmarshal([]byte(m.Data), &message); err != nil {
					log.Printf("fail copying message to model: %s", err)
				} else {
					if message.ID == 1 {
						ms.Start()
						ms.SetFirstReceiving(*message.SentAt)
					}
					ms.AddMeasurement(time.Since(*message.SentAt))
					if message.ID == expectedMessages {
						ms.RegisterLastReceiving()
						last <- true
					}
				}
			}
		}
	}(last, size)

	sub, err := s.client.Subscribe(s.topic, cb)
	if err != nil {
		log.Fatalf("error subscribing: %s", err)
	}
	ready <- true

	log.Println("listening")
	<-last

	defer ms.Stop()
	return sub.Unsubscribe()
}
