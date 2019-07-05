package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
)

type mockData struct {
	EventType string `json:"eventType"`
	Data      string `json:"data"`
}

// ROUTING_KEY=asia.thailand go run producer.go
// ROUTING_KEY=asia.singapore go run producer.go

func main() {
	conn, err := rabbitmq.NewConnection(rabbitmq.ConfigConnection{
		Type: "standalone",
		Url:  "amqp://admin:1234@localhost:5672/",
	})
	if err != nil {
		log.Fatalf("[main]: unable to connect RabbitMQ %+v", err)
	}

	rbMqConfig := rabbitmq.ConfigProducer{
		Exchange: rabbitmq.ConfigExchange{
			Type:       rabbitmq.ExchangeDirect,
			Name:       "asia.exchange",
			RoutingKey: os.Getenv("ROUTING_KEY"),
		},
	}

	producer, err := rabbitmq.NewProducer(conn, rbMqConfig)
	if err != nil {
		log.Fatalf("[main]: unable to initial producer RabbitMQ %+v", err)
	}

	for i := 0; i < 10000; i++ {
		body, _ := json.Marshal(&mockData{
			EventType: "asiaEvent",
			Data:      fmt.Sprintf(os.Getenv("ROUTING_KEY")),
		})

		err = producer.Publish(body)
		if err != nil {
			log.Printf("unable to send message %+v", err)
		}
	}
}
