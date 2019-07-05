package worker

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func SetupExchange(conn *amqp.Connection) {
	// create channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to create channel %+v", err)
	}

	// work exchange
	config := rabbitmq.ConfigExchange{
		Type: rabbitmq.ExchangeTopic,
		Name: "work.exchange",
	}
	err = rabbitmq.ExchangeDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create exchange %+v", err)
	}

	// retry exchange
	config = rabbitmq.ConfigExchange{
		Type: rabbitmq.ExchangeTopic,
		Name: "retry.exchange",
	}
	err = rabbitmq.ExchangeDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create exchange %+v", err)
	}
}
