package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Telefonica/redis-vs-nats/measurement"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisMessassingService struct {
	client *redis.Client
	topic  string
}

func NewRedisMessagingService() (MessagingService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	return &redisMessassingService{client: client, topic: "message"}, nil
}

func (s *redisMessassingService) Publish(message *model.Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("fail marshaling message: %s", err)
	}
	return s.client.Publish(s.topic, string(messageJSON)).Err()
}

func (s *redisMessassingService) Subscribe(ready chan bool, size uint, ms *measurement.Measurements) error {
	pubsub := s.client.Subscribe(s.topic)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		return errors.Wrap(err, "fail subscribing")
	}

	channel := pubsub.ChannelSize(int(size))

	ready <- true
	log.Println("listening")
	for packet := range channel {
		var message model.Message
		if err = json.Unmarshal([]byte(packet.Payload), &message); err != nil {
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
