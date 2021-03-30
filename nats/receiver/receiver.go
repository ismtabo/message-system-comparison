package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/Telefonica/redis-vs-nats/model"

	"github.com/nats-io/nats.go"
)

var wg sync.WaitGroup // 1
var client *nats.Conn
var err error

func init() {
	client, err = nats.Connect("nats://0.0.0.0:4222")
	checkErr(err)
}

func main() {
	wg.Add(1)
	go worker()
	wg.Wait()
}

func worker() {
	// defer wg.Done()

	client.Subscribe("message", func(m *nats.Msg) {
		message := model.Message{}
		json.Unmarshal([]byte(m.Data), &message)
		log.Default().Println("Received message: ", message)
	})

}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
