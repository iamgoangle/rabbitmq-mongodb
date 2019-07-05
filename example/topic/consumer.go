package main

import (
	"log"
	"os"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
)

// ROUTING_KEY=asia.thailand.* QUEUE=asia.thailand.all.province go run consumer.go
// ROUTING_KEY=asia.thailand.# QUEUE=asia.thailand.all.province.road go run consumer.go

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
			Type:    rabbitmq.ExchangeTopic,
			Name:    "asia.exchange.topic",
			Durable: true,
		},
		Queue: rabbitmq.ConfigQueue{
			Name:      os.Getenv("QUEUE"),
			Exclusive: true,
			Bind: rabbitmq.ConfigQueueBind{
				ExchangeName: "asia.exchange.topic",
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
