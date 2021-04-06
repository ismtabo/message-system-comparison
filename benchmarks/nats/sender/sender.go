package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Telefonica/redis-vs-nats/measurement"
	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/nats-io/nats.go"
)

var (
	filename string = "../../../json/100k.json"
	client   *nats.Conn
	ms       measurement.Measurements
	err      error
)

func init() {
	client, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("fail connecting to nats: %s", err)
	}
}

func main() {
	defer client.Close()

	start := time.Now()

	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("fail opening file: %s", err)
	}

	decoder := json.NewDecoder(jsonFile)

	// Read open file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail decoding first token: %s", err)
	}

	for decoder.More() {
		var message model.Message
		err := decoder.Decode(&message)
		if err != nil {
			log.Fatalf("fail decoding message: %s", err)
		}

		now := time.Now()
		message.SentAt = &now

		data, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("fail marshaling message: %s", err)
		}

		start := time.Now()
		if err := client.Publish("message", data); err != nil {
			log.Fatalf("fail sending message: %s", err)
		}
		ms.AddMeasurement(time.Since(start))

	}

	// Close the file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail decoding last token: %s", err)
	}

	elapsed := time.Since(start)
	log.Printf("Nats Sender took %s\n", elapsed)
	log.Printf("Latency: mean %gs, variance %gs Throughput: %g msg/s", ms.Mean(), ms.Var(), 1/ms.Mean())
	log.Println("Throughput:", 1/ms.Mean(), "msg/s")

}
