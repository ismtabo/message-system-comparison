package benchmark

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Telefonica/redis-vs-nats/measurement"
	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/Telefonica/redis-vs-nats/services"
)

type Benchmark struct {
	service  services.MessagingService
	filename string
	size     uint
}

func NewBenchmark(messaging services.MessagingService) *Benchmark {
	filename := "../../json/100k.json"
	var size uint = 100_000
	return &Benchmark{service: messaging, filename: filename, size: size}
}

func (b *Benchmark) Run() {
	ready := make(chan bool)
	done := make(chan bool)

	var receiverMs measurement.Measurements
	go b.runReceiver(&receiverMs, ready, done)
	<-ready

	var senderMs measurement.Measurements
	b.runSender(&senderMs)
	<-done

	log.Printf("Sender took: %.5g s", senderMs.Elapsed().Seconds())
	log.Printf("Messages sent: %d", senderMs.N())
	log.Printf("Latency: mean %.5g s, variance %.5g s", senderMs.Mean(), senderMs.Var())
	log.Printf("Throughput: %.5g msg/s", 1/senderMs.Mean())
	log.Println("-------------------")
	log.Printf("Receiver took: %.5g s", receiverMs.Elapsed().Seconds())
	log.Printf("Time among first send and last receive: %.5g s", receiverMs.ElapsedMessagingTime().Seconds())
	log.Printf("Messages received: %d", receiverMs.N())
	log.Printf("Transmission time: mean %.5g s, variance %.5g s", receiverMs.Mean(), receiverMs.Var())
	log.Printf("Bandwidth: %.5g msg/s", 1/receiverMs.Mean())
}

func (b *Benchmark) runSender(ms *measurement.Measurements) {

	jsonFile, err := os.Open(b.filename)
	if err != nil {
		log.Fatalf("fail opening file: %s", err)
	}

	decoder := json.NewDecoder(jsonFile)

	// Read open file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail decoding first token: %s", err)
	}

	ms.Start()
	for decoder.More() {
		var message model.Message
		err := decoder.Decode(&message)
		if err != nil {
			log.Fatalf("fail decoding message: %s", err)
		}

		now := time.Now()
		message.SentAt = &now

		start := time.Now()
		if err := b.service.Publish(&message); err != nil {
			log.Fatalf("fail sending message: %s", err)
		}
		ms.AddMeasurement(time.Since(start))

	}
	ms.Stop()

	// Close the file
	_, err = decoder.Token()
	if err != nil {
		log.Fatalf("fail decoding last token: %s", err)
	}
}

func (b *Benchmark) runReceiver(ms *measurement.Measurements, ready chan bool, done chan bool) {
	defer close(done)
	if err := b.service.Subscribe(ready, b.size, ms); err != nil {
		log.Fatalf("fail subscribing: %s", err)
	}
}
