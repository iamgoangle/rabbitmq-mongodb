package rabbitmq

import (
	"errors"

	"log"

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

	log.Println("[QueueDeclare]: new queue %v", c.Name)
	return q, nil
}

// QueueBind binding my queue with exchange routing
func QueueBind(config ConfigQueue, ch *amqp.Channel) error {
	log.Println("config: %+v", config.Bind)

	if err := ch.QueueBind(
		config.Name,
		config.Bind.RoutingKey,
		config.Bind.ExchangeName,
		config.NoWait,
		nil,
	); err != nil {
		return errors.New("[QueueBind]: unable to queue bind" + err.Error())
	}

	return nil
}

// IsQueueExist checking my queue is already declare in RabbitMQ
func IsQueueExist(name string, ch *amqp.Channel) bool {
	var exist bool
	_, err := ch.QueueInspect(name)
	if err == nil {
		exist = true
	}

	return exist
}
