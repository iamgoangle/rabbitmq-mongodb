package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
)

type mockData struct {
	EventType string `json:"eventType"`
	Data      string `json:"data"`
}

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
			Name:       "work.exchange",
			RoutingKey: "work.routing",
		},
	}

	producer, err := rabbitmq.NewProducer(conn, rbMqConfig)
	if err != nil {
		log.Fatalf("[main]: unable to initial producer RabbitMQ %+v", err)
	}

	for i := 0; i < 10; i++ {
		body, _ := json.Marshal(&mockData{
			EventType: "newEvent",
			Data:      fmt.Sprintf("event id = %+v", i),
		})

		err = producer.Publish(body)
		if err != nil {
			log.Printf("unable to send my event %+v", err)
		}

		time.Sleep(time.Second)
	}
}
