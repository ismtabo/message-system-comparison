package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Telefonica/redis-vs-nats/benchmarks/kafka-streams/codec"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/jinzhu/copier"
	"github.com/lovoo/goka"
)

var (
	brokers               = []string{"localhost:9092"}
	topic     goka.Stream = "messages-stream"
	group     goka.Group  = "messages-group"
	processor *goka.Processor
	size      int64 = 100_000
	err       error
	n         int64
	Ex, Ex2   float64
)

func main() {
	last := make(chan bool)

	// process callback is invoked for each message delivered from
	// "example-stream" topic.
	cb := func(last chan bool, expectedMessages int64) func(ctx goka.Context, msg interface{}) {
		return func(ctx goka.Context, msg interface{}) {
			if n < int64(expectedMessages) {
				var message model.Message
				if err := copier.Copy(&message, msg); err != nil {
					log.Printf("fail copying message to model: %s", err)
				} else {
					x := time.Since(*message.SentAt).Seconds()
					n++
					Ex += x
					Ex2 += x * x
					if n >= int64(expectedMessages) {
						last <- true
					}
				}
			}
		}
	}(last, size)

	// Define a new processor group. The group defines all inputs, outputs, and
	// serialization formats. The group-table topic is "example-group-table".
	g := goka.DefineGroup(group,
		goka.Input(topic, new(codec.MessageCodec), cb),
	)

	processor, err = goka.NewProcessor(brokers, g)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
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

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-wait: // wait for SIGINT/SIGTERM
		log.Println("canceling")
	case <-last:
		log.Println("retrieve last")
	}
	cancel() // gracefully stop processor
	<-done

	fN := float64(n)
	fEx := float64(Ex)
	fEx2 := float64(Ex2)

	mean := fEx / fN
	variance := (fEx2 - (fEx*fEx)/fN) / (fN - 1)
	log.Printf("Messages received: %d Latency: mean %gs, variance %gs Throughput: %g msg/s", n, mean, variance, 1/mean)
}
