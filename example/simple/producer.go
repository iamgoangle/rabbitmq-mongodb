package main

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
)

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
			Type: "",
		},
		Queue: rabbitmq.ConfigQueue{
			Name: "hello-simple",
		},
	}

	producer, err := rabbitmq.NewProducer(conn, rbMqConfig)
	if err != nil {
		log.Fatalf("[main]: unable to connect RabbitMQ %+v", err)
	}

	for i := 0; i < 10; i++ {
		body := []byte(`{"eventType":"sendMessage","data":"Hello World","Time":1294706395881547000}`)
		err = producer.Publish(body)
		if err != nil {
			log.Printf("unable to send my message %+v", err)
		}
	}
}
