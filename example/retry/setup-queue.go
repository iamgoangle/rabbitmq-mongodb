package main

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
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

	// ===========================
	// work queue
	// ===========================
	config := rabbitmq.ConfigQueue{
		Name: "work.queue",
		Bind: rabbitmq.ConfigQueueBind{
			RoutingKey:   "work.routing",
			ExchangeName: "work.exchange",
		},
	}
	q, err := rabbitmq.QueueDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create queue %+v", err)
	}

	// bind q with exchange
	err = rabbitmq.QueueBind(config, ch, q)
	if err != nil {
		log.Fatalf("unable to binding queue %+v", err)
	}

	// ===========================
	// retry queue
	// ===========================
	argsRetryExchange := make(amqp.Table)
	argsRetryExchange["x-dead-letter-exchange"] = "work.exchange"
	argsRetryExchange["x-dead-letter-routing-key"] = "work.routing"
	argsRetryExchange["x-message-ttl"] = 10000
	config = rabbitmq.ConfigQueue{
		Name: "retry.queue",
		Bind: rabbitmq.ConfigQueueBind{
			RoutingKey:   "retry.routing",
			ExchangeName: "retry.exchange",
		},
		Args: argsRetryExchange,
	}
	q, err = rabbitmq.QueueDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create queue %+v", err)
	}

	// bind q with exchange
	err = rabbitmq.QueueBind(config, ch, q)
	if err != nil {
		log.Fatalf("unable to binding queue %+v", err)
	}
}
