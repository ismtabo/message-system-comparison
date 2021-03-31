package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/nats-io/nats.go"
)

var (
	client  *nats.Conn
	size    int64 = 100_000
	err     error
	n       int64
	Ex, Ex2 float64
)

func init() {
	client, err = nats.Connect("nats://0.0.0.0:4222")
	if err != nil {
		log.Fatalf("error connecting to nats: %s", err)
	}
}

func main() {
	last := make(chan bool)
	cb := func(last chan bool, expectedMessages int64) func(*nats.Msg) {
		return func(m *nats.Msg) {
			if n < expectedMessages {
				var message model.Message
				if err := json.Unmarshal([]byte(m.Data), &message); err != nil {
					log.Printf("fail copying message to model: %s", err)
				} else {
					x := time.Since(*message.SentAt).Seconds()
					n++
					Ex += x
					Ex2 += x * x
					if n >= expectedMessages {
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

	fN := float64(n)
	fEx := float64(Ex)
	fEx2 := float64(Ex2)

	mean := fEx / fN
	variance := (fEx2 - (fEx*fEx)/fN) / (fN - 1)
	log.Printf("Messages received: %d Latency: mean %gs, variance %gs Throughput: %g msg/s", n, mean, variance, 1/mean)
}
