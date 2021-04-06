package main

import (
	"log"

	"github.com/Telefonica/redis-vs-nats/benchmark"
	"github.com/Telefonica/redis-vs-nats/services"
)

func main() {
	messagingServices := map[string]func() (services.MessagingService, error){
		"Kafka":    services.NewKafkaMessagingService,
		"Nats":     services.NewNatsMessagingService,
		"RabbitMq": services.NewRabbitMqMessagingService,
		"Redis":    services.NewRedisMessagingService,
	}

	for title, messagingServiceFn := range messagingServices {
		messagingService, err := messagingServiceFn()
		if err != nil {
			log.Fatal(err)
		}
		b := benchmark.NewBenchmark(messagingService)
		log.Printf("****Running %s benchmark****", title)
		b.Run()
		log.Println("****************************")
	}
}
