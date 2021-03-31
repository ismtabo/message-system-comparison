package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/go-redis/redis"
)

var (
	client  *redis.Client
	size    int64 = 100_000
	n       int64
	Ex, Ex2 float64
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})
}

func main() {
	pubsub := client.Subscribe("message")
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	channel := pubsub.ChannelSize(int(size))

	log.Println("listening")
	for packet := range channel {
		var message model.Message
		if err = json.Unmarshal([]byte(packet.Payload), &message); err != nil {
			log.Printf("error unmarshaling message: %s", err)
		} else {
			x := time.Since(*message.SentAt).Seconds()
			n++
			Ex += x
			Ex2 += x * x
			if n >= size {
				break
			}
		}
	}

	fN := float64(n)
	fEx := float64(Ex)
	fEx2 := float64(Ex2)

	mean := fEx / fN
	variance := (fEx2 - (fEx*fEx)/fN) / (fN - 1)
	log.Printf("Messages received: %d Latency: mean %gs, variance %gs Throughput: %g msg/s\n", n, mean, variance, 1/mean)
}
