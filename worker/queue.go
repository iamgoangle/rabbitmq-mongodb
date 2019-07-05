package worker

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func SetupQueue(conn *amqp.Connection) {
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
			RoutingKey:   "work.routing.#",
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
	argsRetryExchange["x-dead-letter-routing-key"] = "work.routing.retry"
	argsRetryExchange["x-message-ttl"] = 5000
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
