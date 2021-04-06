package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/Telefonica/redis-vs-nats/measurement"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/Telefonica/redis-vs-nats/services/dto"
	"github.com/lovoo/goka"
)

type kafkaMessagingService struct {
	emitter *goka.Emitter
	brokers []string
	topic   goka.Stream
	group   goka.Group
}

func NewKafkaMessagingService() (MessagingService, error) {
	var brokers = []string{"localhost:9092"}
	var topic goka.Stream = "messages-stream"
	var group goka.Group = "messages-group"

	emitter, err := goka.NewEmitter(brokers, topic, new(dto.MessageCodec))
	if err != nil {
		return nil, errors.Wrap(err, "fail creating emitter")
	}

	return &kafkaMessagingService{emitter: emitter, brokers: brokers, topic: topic, group: group}, nil
}

func (s *kafkaMessagingService) Publish(message *model.Message) error {
	_, err := s.emitter.Emit(fmt.Sprint(message.ID), message)
	return err
}

func (s *kafkaMessagingService) Subscribe(ready chan bool, size uint, ms *measurement.Measurements) error {
	last := make(chan bool)

	// process callback is invoked for each message delivered from
	// "example-stream" topic.
	cb := func(last chan bool, expectedMessages uint) func(ctx goka.Context, msg interface{}) {
		return func(ctx goka.Context, msg interface{}) {
			if ms.N() < expectedMessages {
				if message, ok := msg.(*model.Message); !ok {
					log.Printf("fail casting received message: %+v", msg)
				} else {
					if message.ID == 1 {
						ms.Start()
						ms.SetFirstReceiving(*message.SentAt)
					}
					x := time.Since(*message.SentAt)
					ms.AddMeasurement(x)
					if message.ID == expectedMessages {
						ms.RegisterLastReceiving()
						last <- true
					}
				}
			}
		}
	}(last, size)

	// Define a new processor group. The group defines all inputs, outputs, and
	// serialization formats. The group-table topic is "example-group-table".
	g := goka.DefineGroup(s.group,
		goka.Input(s.topic, new(dto.MessageCodec), cb),
	)

	processor, err := goka.NewProcessor(s.brokers, g)
	if err != nil {
		return errors.Wrap(err, "error creating processor")
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go func() {
		defer close(done)
		log.Println("listening")
		if err = processor.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()
	ready <- true

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-wait: // wait for SIGINT/SIGTERM
		log.Println("canceling")
	case <-last:
		log.Println("retreive last")
	}
	cancel() // gracefully stop processor
	<-done

	ms.Stop()

	return nil
}
