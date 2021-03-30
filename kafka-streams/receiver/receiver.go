package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Telefonica/redis-vs-nats/kafka-streams/codec"
	"github.com/lovoo/goka"
)

var (
	brokers             = []string{"localhost:9092"}
	topic   goka.Stream = "example-stream"
	group   goka.Group  = "example-group"
)
var processor *goka.Processor
var err error

func init() {
	// process callback is invoked for each message delivered from
	// "example-stream" topic.
	cb := func(ctx goka.Context, msg interface{}) {
		log.Printf("key = %s, msg = %v", ctx.Key(), msg)
	}
	// Define a new processor group. The group defines all inputs, outputs, and
	// serialization formats. The group-table topic is "example-group-table".
	g := goka.DefineGroup(group,
		goka.Input(topic, new(codec.MessageCodec), cb),
	)

	processor, err = goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go func() {
		defer close(done)
		if err = processor.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait   // wait for SIGINT/SIGTERM
	cancel() // gracefully stop processor
	<-done
}
