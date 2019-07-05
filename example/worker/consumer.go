package main

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
)

func main() {
	rbMqConfig := rabbitmq.ConfigConsumer{
		Conn: rabbitmq.ConfigConnection{
			Type: "standalone",
			Url:  "amqp://admin:1234@localhost:5672/",
		},
		// direct type
		// Exchange: rabbitmq.ConfigExchange{
		// 	Type: rabbitmq.ExchangeDirect,
		// 	Name: "Campaign-test",
		// },

		// topic type
		Exchange: rabbitmq.ConfigExchange{
			Type:       rabbitmq.ExchangeTopic,
			Name:       "Campaign-test-topic",
			RoutingKey: "campaign.pool.*",
		},
		Queue: rabbitmq.ConfigQueue{
			Name:            "hello-direct",
			DeleteWhenUnuse: false,
			Durable:         true,
			Exclusive:       false,
			NoWait:          false,
		},
	}

	consumer, err := rabbitmq.NewConsumer(rbMqConfig)
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
