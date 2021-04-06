package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/nats-io/nats.go"
)

var (
	client            *nats.Conn
	size              uint = 100_000
	err               error
	n                 uint
	Ex, Ex2           float64
	start             time.Time
	firstTime         time.Time
	sinceFirstMessage time.Duration
)

func init() {
	client, err = nats.Connect("nats://0.0.0.0:4222")
	if err != nil {
		log.Fatalf("error connecting to nats: %s", err)
	}
}

func main() {
	last := make(chan bool)
	cb := func(last chan bool, expectedMessages uint) func(*nats.Msg) {
		return func(m *nats.Msg) {
			if n < expectedMessages {
				var message model.Message
				if err := json.Unmarshal([]byte(m.Data), &message); err != nil {
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

	sub, err := client.Subscribe("message", cb)
	if err != nil {
		log.Fatalf("error subscribing: %s", err)
	}

	log.Println("listening")
	<-last
	sub.Unsubscribe()
	if err != nil {
		log.Fatalf("error unsubscribing: %s", err)
	}

	elapsed := time.Since(start)

	fN := float64(n)
	fEx := float64(Ex)
	fEx2 := float64(Ex2)

	mean := fEx / fN
	variance := (fEx2 - (fEx*fEx)/fN) / (fN - 1)
	log.Printf("Messages received: %d Total time: %s Time among first send and last receive: %s", n, elapsed, sinceFirstMessage)
	log.Printf("Performance: mean %gs, variance %gs", mean, variance)
	log.Println("Throughput:", 1/mean, "msg/s")
}
