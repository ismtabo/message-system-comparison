package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/nats-io/nats.go"
)

var wg sync.WaitGroup // 1
var client *nats.Conn
var err error
var n int64
var Ex, Ex2 float64

func init() {
	client, err = nats.Connect("nats://0.0.0.0:4222")
	checkErr(err)
}

func main() {
	wg.Add(1)
	go worker()
	wg.Wait()

	fN := float64(n)
	fEx := float64(Ex)
	fEx2 := float64(Ex2)

	mean := fEx / fN
	variance := (fEx2 - (fEx*fEx)/fN) / (fN - 1)
	log.Printf("Messages received: %d\nLatency: mean %gs, variance %gs\n", n, mean, variance)
}

func worker() {
	defer wg.Done()

	last := make(chan bool)
	cb := func(last chan bool) func(*nats.Msg) {
		return func(m *nats.Msg) {
			if n < 100_000 {
				var message model.Message
				if err := json.Unmarshal([]byte(m.Data), &message); err != nil {
					log.Printf("fail copying message to model: %s", err)
				} else {
					x := time.Since(*message.SentAt).Seconds()
					n++
					Ex += x
					Ex2 += x * x
					if n >= 100_000 {
						last <- true
					}
				}
			}
		}
	}(last)

	sub, err := client.Subscribe("message", cb)
	if err != nil {
		log.Fatalf("error subscribing: %s", err)
	}

	<-last
	sub.Unsubscribe()
	if err != nil {
		log.Fatalf("error unsubscribing: %s", err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
