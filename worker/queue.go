package worker

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
)

// InitQueue work and delay queue
func InitQueue(conn *amqp.Connection) {
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
	_, err = rabbitmq.QueueDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create queue %+v", err)
	}

	// bind q with exchange
	err = rabbitmq.QueueBind(config, ch)
	if err != nil {
		log.Fatalf("unable to binding queue %+v", err)
	}

	// ===========================
	// delay queue
	// ===========================
	argsDelayQueue := make(amqp.Table)
	argsDelayQueue["x-dead-letter-exchange"] = "work.exchange"
	argsDelayQueue["x-dead-letter-routing-key"] = "work.routing.delay"
	argsDelayQueue["x-message-ttl"] = 5000
	config = rabbitmq.ConfigQueue{
		Name: "delay.queue",
		Bind: rabbitmq.ConfigQueueBind{
			RoutingKey:   "delay.routing",
			ExchangeName: "delay.exchange",
		},
		Args: argsDelayQueue,
	}
	_, err = rabbitmq.QueueDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create queue %+v", err)
	}

	// bind q with exchange
	err = rabbitmq.QueueBind(config, ch)
	if err != nil {
		log.Fatalf("unable to binding queue %+v", err)
	}
}
