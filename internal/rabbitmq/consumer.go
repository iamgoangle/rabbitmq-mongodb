package rabbitmq

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
)

type Consumer interface {
	// WorkerProcessor spawn queue consume worker
	WorkerProcessor() (deliver <-chan amqp.Delivery, err error)

	// Close channels and connection
	Close()
}

type ConfigConsumer struct {
	Exchange ConfigExchange
	Queue    ConfigQueue
}

type consumer struct {
	ch     *amqp.Channel
	config ConfigConsumer
	conn   *amqp.Connection
	// q      amqp.Queue
}

const (
	prefetch       = 100
	prefetchSize   = 0
	prefetchGlobal = false

	autoAck     = false
	consumerTag = ""
	exclusive   = false
	noLocal     = false
	noWait      = false
)

// NewConsumer initialise new consumer
func NewConsumer(conn *amqp.Connection, config ConfigConsumer) (Consumer, error) {
	ch, err := conn.Channel()
	failOnError(err, "[NewConsumer]: Failed to open a channel")

	err = ExchangeDeclare(config.Exchange, ch)
	if err != nil {
		return nil, err
	}

	q, err := QueueDeclare(config.Queue, ch)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(
		prefetch,
		prefetchSize,
		prefetchGlobal,
	)
	if err != nil {
		return nil, errors.New("[WorkerProcessor]: unable to create QoS")
	}

	err = QueueBind(config.Queue, ch, q)
	if err != nil {
		return nil, err
	}

	return &consumer{
		conn:   conn,
		config: config,
		ch:     ch,
	}, nil
}

// WorkerProcessor spawn queue consume worker
func (c *consumer) WorkerProcessor() (deliver <-chan amqp.Delivery, err error) {
	log.Println("[WorkerProcessor]: start worker")

	deliver, err = c.ch.Consume(
		// c.q.Name,
		c.config.Queue.Name,
		consumerTag,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		nil,
	)

	return
}

func (c *consumer) QueueBind() (err error) {
	err = c.ch.QueueBind(
		c.config.Queue.Name,            // queue name
		c.config.Queue.Bind.RoutingKey, // routing key
		c.config.Exchange.Name,         // exchange
		c.config.Queue.NoWait,
		nil)

	log.Println("[QueueBind]: " + c.config.Queue.Bind.RoutingKey)

	return
}

// Close channels and connection
func (c *consumer) Close() {
	c.conn.Close()
	c.ch.Close()
}
