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
			Type:       rabbitmq.ExchangeFanout,
			Name:       "asia.exchange.fanout",
			Durable:    true,
			RoutingKey: "",
		},
	}

	producer, err := rabbitmq.NewProducer(conn, rbMqConfig)
	if err != nil {
		log.Fatalf("[main]: unable to initial producer RabbitMQ %+v", err)
	}

	for i := 0; i < 10000; i++ {
		body, _ := json.Marshal(&mockData{
			EventType: "asiaEvent",
			Data:      fmt.Sprintf("test fanout %+v", i),
		})

		err = producer.Publish(body)
		if err != nil {
			log.Printf("unable to send message %+v", err)
		}

		time.Sleep(1 * time.Second)
	}
}
