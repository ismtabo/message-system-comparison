package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Telefonica/redis-vs-nats/kafka-streams/codec"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/lovoo/goka"
)

var (
	brokers              = []string{"localhost:9092"}
	topic    goka.Stream = "example-stream"
	filename string      = "../../json/100k.json"
	emitter  *goka.Emitter
	err      error
)

func init() {
	emitter, err = goka.NewEmitter(brokers, topic, new(codec.MessageCodec))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
}

func main() {
	defer emitter.Finish()

	start := time.Now()

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error decoding message: %v", err)
	}

	decoder := json.NewDecoder(jsonFile)

	// Read opening file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("error decoding token: %v", err)
	}

	var message model.Message
	for decoder.More() {
		err := decoder.Decode(&message)
		if err != nil {
			log.Fatalf("error decoding message: %v", err)
		}
		now := time.Now()
		message.SentAt = &now
		_, err = emitter.Emit(fmt.Sprint(message.ID), message)
		if err != nil {
			log.Fatalf("error emitting message: %v", err)
		}
	}

	elapsed := time.Since(start)
	log.Printf("Kafka Streams Sender took %s\n", elapsed)
}
