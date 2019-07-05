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
		log.Fatal("[main]: unable to connect RabbitMQ %+v", err)
	}

	// create channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to create channel %+v", err)
	}

	// work exchange
	config := rabbitmq.ConfigExchange{
		Type:       rabbitmq.ExchangeDirect,
		Name:       "work.exchange",
		RoutingKey: "work.routing",
	}
	err = rabbitmq.ExchangeDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create exchange %+v", err)
	}

	// retry exchange
	config = rabbitmq.ConfigExchange{
		Type:       rabbitmq.ExchangeDirect,
		Name:       "retry.exchange",
		RoutingKey: "retry.routing",
	}
	err = rabbitmq.ExchangeDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create exchange %+v", err)
	}
}
