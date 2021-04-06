package services

import (
	"github.com/Telefonica/redis-vs-nats/measurement"
	"github.com/Telefonica/redis-vs-nats/model"
)

type MessagingService interface {
	Publish(*model.Message) error
	Subscribe(chan bool, uint, *measurement.Measurements) error
}
