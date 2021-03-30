package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/go-redis/redis"
)

var wg sync.WaitGroup // 1
var client *redis.Client
var n int64
var Ex, Ex2 float64

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})
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

	pubsub := client.Subscribe("message")
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	channel := pubsub.ChannelSize(1000000)

	for packet := range channel {
		var message model.Message
		if err = json.Unmarshal([]byte(packet.Payload), &message); err != nil {
			log.Printf("error unmarshaling message: %s", err)
		} else {
			x := time.Since(*message.SentAt).Seconds()
			n++
			Ex += x
			Ex2 += x * x
			if n >= 100_000 {
				break
			}
		}
	}
}
