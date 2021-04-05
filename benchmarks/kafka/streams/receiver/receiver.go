package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Telefonica/redis-vs-nats/benchmarks/kafka/streams/codec"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/lovoo/goka"
)

var (
	brokers                       = []string{"localhost:9092"}
	topic             goka.Stream = "messages-stream"
	group             goka.Group  = "messages-group"
	processor         *goka.Processor
	size              uint = 100_000
	err               error
	n                 int64
	Ex, Ex2           float64
	start             time.Time
	firstTime         time.Time
	sinceFirstMessage time.Duration
)

func main() {
	last := make(chan bool)

	// process callback is invoked for each message delivered from
	// "example-stream" topic.
	cb := func(last chan bool, expectedMessages uint) func(ctx goka.Context, msg interface{}) {
		return func(ctx goka.Context, msg interface{}) {
			if n < int64(expectedMessages) {
				if message, ok := msg.(*model.Message); !ok {
					log.Printf("fail copying message to model: %s", err)
				} else {
					if message.ID == 1 {
						start = time.Now()
						firstTime = *message.SentAt
					}
					x := time.Since(*message.SentAt).Seconds()
					n++
					Ex += x
					Ex2 += x * x
					if message.ID == expectedMessages {
						sinceFirstMessage = time.Since(firstTime)
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
		log.Println("retreive last")
	}
	cancel() // gracefully stop processor
	<-done

	elapsed := time.Since(start)

	fN := float64(n)
	fEx := float64(Ex)
	fEx2 := float64(Ex2)

	mean := fEx / fN
	variance := (fEx2 - (fEx*fEx)/fN) / (fN - 1)
	log.Println("fN", fmt.Sprintf("%g", fN), "fEx", fmt.Sprintf("%g", fEx))
	log.Printf("Messages received: %d Total time: %s Time among first send and last receive: %s", n, elapsed, sinceFirstMessage)
	log.Printf("Latency: mean %gs, variance %gs Throughput: %g msg/s", mean, variance, 1/mean)
}
