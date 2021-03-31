package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/go-redis/redis"
)

var (
	filename string = "../../../json/100k.json"
	client   *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
}

func main() {

	start := time.Now()

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("fail open json file: %s", err)
	}

	decoder := json.NewDecoder(jsonFile)

	// Read opening file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail decoding first token: %s", err)
	}

	var message model.Message
	for decoder.More() {
		err := decoder.Decode(&message)
		if err != nil {
			log.Fatalf("fail decoding message: %s", err)
		}

		now := time.Now()
		message.SentAt = &now

		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("fail marshaling message: %s", err)
		}

		_, err = client.Publish("message", string(messageJSON)).Result()
		if err != nil {
			log.Fatalf("fail sending message '%d': %s", message.ID, err)
		}
	}

	// Close the file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail decoding last token: %s", err)
	}

	elapsed := time.Since(start)
	log.Printf("Redis Sender took %s\n", elapsed)

}
