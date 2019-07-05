package rabbitmq

import (
	"errors"
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/empty"
	amqp "github.com/streadway/amqp"
)

// Config queue property
type ConfigQueue struct {
	Name            string
	Durable         bool
	DeleteWhenUnuse bool
	Exclusive       bool
	NoWait          bool
	Bind            ConfigQueueBind
	Args            amqp.Table
}

type ConfigQueueBind struct {
	// Name of queue
	Name string

	RoutingKey string

	ExchangeName string
}

// QueueDeclare create a queue meta data
func QueueDeclare(c ConfigQueue, ch *amqp.Channel) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		c.Name,
		c.Durable,
		c.DeleteWhenUnuse,
		c.Exclusive,
		c.NoWait,
		c.Args,
	)
	if err != nil {
		return amqp.Queue{},
			errors.New("[QueueDeclare]: unable to create queue " + err.Error())
	}

	log.Printf("[QueueDeclare]: new queue %v", c.Name)
	return q, nil
}

// QueueBind binding my queue with exchange routing
func QueueBind(config ConfigQueue, ch *amqp.Channel, q amqp.Queue) error {
	if !empty.StructIsEmpty(config.Bind.ExchangeName) {
		if err := ch.QueueBind(
			config.Name,
			config.Bind.RoutingKey,
			config.Bind.ExchangeName,
			config.NoWait,
			nil,
		); err != nil {
			return errors.New("[QueueBind]: unable to queue bind" + err.Error())
		}
	}
	log.Println("[QueueBind]: " + config.Bind.RoutingKey)

	return nil
}
