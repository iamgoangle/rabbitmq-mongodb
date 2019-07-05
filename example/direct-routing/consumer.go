package main

import (
	"log"
	"os"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
)

// ROUTING_KEY=asia.thailand QUEUE=asia.thailand.queue go run consumer.go
// ROUTING_KEY=asia.singapore QUEUE=asia.singapore.queue go run consumer.go
func main() {
	conn, err := rabbitmq.NewConnection(rabbitmq.ConfigConnection{
		Type: "standalone",
		Url:  "amqp://admin:1234@localhost:5672/",
	})
	if err != nil {
		log.Fatalln("[main]: unable to connect RabbitMQ %+v", err)
	}

	rbMqConfig := rabbitmq.ConfigConsumer{
		Exchange: rabbitmq.ConfigExchange{
			Type: rabbitmq.ExchangeDirect,
			Name: "asia.exchange",
		},
		Queue: rabbitmq.ConfigQueue{
			Name: os.Getenv("QUEUE"),
			Bind: rabbitmq.ConfigQueueBind{
				ExchangeName: "asia.exchange",
				RoutingKey:   os.Getenv("ROUTING_KEY"),
			},
		},
	}

	consumer, err := rabbitmq.NewConsumer(conn, rbMqConfig)
	if err != nil {
		log.Panic(err)
	}

	msgs, err := consumer.WorkerProcessor()
	if err != nil {
		log.Panic(err)
	}

	forever := make(chan bool)

	go func() {
		for m := range msgs {
			log.Println("Received a message: %s", string(m.Body))
			log.Println("Done")
			m.Ack(false)
		}
	}()

	<-forever

	consumer.Close()
}
