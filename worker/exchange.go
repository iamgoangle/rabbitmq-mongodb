package worker

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
)

// InitExchange rabbitmq routing
func InitExchange(conn *amqp.Connection) {
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

	// delay exchange
	config = rabbitmq.ConfigExchange{
		Type: rabbitmq.ExchangeTopic,
		Name: "delay.exchange",
	}
	err = rabbitmq.ExchangeDeclare(config, ch)
	if err != nil {
		log.Fatalf("unable to create exchange %+v", err)
	}
}
