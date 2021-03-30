package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Telefonica/redis-vs-nats/model"
	"github.com/bxcodec/faker"
)

func main() {

	datasets := map[string]int{
		"1.json":    1,
		"10.json":   10,
		"1k.json":   1_000,
		"10k.json":  10_000,
		"100k.json": 100_000,
	}

	for filename, value := range datasets {
		message := model.Message{}
		messages := []model.Message{}

		// Added +1 to start with 1 not 0
		for i := 1; i < value+1; i++ {
			if err := faker.FakeData(&message); err != nil {
				log.Fatalf("fail faking data: %s", err)
			}
			message.ID = uint(i)
			now := time.Now()
			message.CreatedAt = &now
			message.UpdatedAt = message.CreatedAt
			messages = append(messages, message)
		}

		f, err := os.Create("../json/" + filename)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		json, _ := json.MarshalIndent(messages, "", "    ")
		err = ioutil.WriteFile("../json/"+filename, json, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

}
