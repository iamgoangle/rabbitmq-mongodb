package rabbitmq

import (
	"errors"
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/empty"
	amqp "github.com/streadway/amqp"
)

type ConfigExchange struct {
	Name            string
	Type            string
	Durable         bool
	DeleteWhenUnuse bool
	Internal        bool
	NoWait          bool

	RoutingKey string
	Args       amqp.Table
}

func exchangeTypeFactory(t string, body []byte) amqp.Publishing {
	switch t {
	case "json":
		return amqp.Publishing{
			ContentType: "text/json",
			Body:        body,
		}
	default:
		return amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		}
	}
}

func ExchangeDeclare(config ConfigExchange, ch *amqp.Channel) error {
	// FIXME: create exchange if it does not exists
	if !empty.StructIsEmpty(config.Type) {
		if err := ch.ExchangeDeclare(
			config.Name,
			config.Type,
			config.Durable,
			config.DeleteWhenUnuse,
			config.Internal,
			config.NoWait,
			config.Args, // arguments
		); err != nil {
			return errors.New("[NewProducer]: unable to create queue exchange " + err.Error())
		}
	}

	log.Printf("[ExchangeDeclare]: exchange create %v", config.Name)

	return nil
}
