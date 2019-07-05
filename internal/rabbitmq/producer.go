package rabbitmq

import (
	"errors"
	"log"

	amqp "github.com/streadway/amqp"
)

const (
	ExchangeDirect  = "direct"
	ExchangeFanout  = "fanout"
	ExchangeTopic   = "topic"
	ExchangeHeaders = "headers"
)

// Producer behaviors
type Producer interface {
	// Publish factory method
	Publish(body []byte) error

	// PublishDirectExchange handles direct with routing key exchange
	PublishDirectExchange(body []byte) error

	// PublishDefaultExchange handles default send to queue
	PublishDefaultExchange(body []byte) error

	// PublishFanoutExchange handles fanout exchange pub/sub
	PublishFanoutExchange(body []byte) error

	// PublishTopicExchange handles matching pattern exchange
	PublishTopicExchange(body []byte) error

	// Release connection
	Close()
}

// Config producer
type ConfigProducer struct {
	Exchange ConfigExchange
	Queue    ConfigQueue
}

// Channel wrapper amqp.Channel
type Channel struct {
	*amqp.Channel
}

// Queue wrapper amqp.Queue
type Queue struct {
	amqp.Queue
}

type producer struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config ConfigProducer
}

// NewProducer amqp producer actor
func NewProducer(conn *amqp.Connection, config ConfigProducer) (Producer, error) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ExchangeDeclare(config.Exchange, ch)
	if err != nil {
		return nil, err
	}

	// _, err = QueueDeclare(config.Queue, ch)
	// if err != nil {
	// 	return nil, err
	// }

	// err = QueueBind(config.Queue, ch, q)
	// if err != nil {
	// 	return nil, err
	// }

	return &producer{
		ch:     ch,
		config: config,
		conn:   conn,
	}, nil
}

// Publish factory method to exchange the message to queue by type
// TODO: To imrove init exchange type when NewProducer by pointer it to exechange type func
func (p *producer) Publish(body []byte) error {
	switch p.config.Exchange.Type {
	case ExchangeDirect:
		return p.PublishDirectExchange(body)
	case ExchangeFanout:
		return p.PublishFanoutExchange(body)
	case ExchangeTopic:
		return p.PublishTopicExchange(body)
	default:
		return p.PublishDefaultExchange(body)
	}
}

// PublishDirectExchange handles default exchange
// Direct: A direct exchange delivers messages to queues based on a message routing key.
// In a direct exchange, the message is routed to the queues whose binding key exactly matches the routing key of the message.
func (p *producer) PublishDirectExchange(body []byte) error {
	err := p.publish(p.config.Exchange.Name, p.config.Exchange.RoutingKey, body)
	if err != nil {
		return errors.New("[PublishDirectExchange]: unable to publish a message " + err.Error())
	}

	log.Println("[PublishDirectExchange]: publish message to queue")

	return nil
}

// PublishDefaultExchange handles default exchange
func (p *producer) PublishDefaultExchange(body []byte) error {
	err := p.publish("", p.config.Queue.Name, body)
	if err != nil {
		return errors.New("[PublishDefaultExchange]: unable to publish a message " + err.Error())
	}

	log.Println("[PublishDefaultExchange]: publish message to queue")

	return nil
}

// PublishFanoutExchange broadcasts all the messages to queue that subscribe my exchange
func (p *producer) PublishFanoutExchange(body []byte) error {
	err := p.publish(p.config.Exchange.Name, "", body)
	if err != nil {
		return errors.New("[PublishFanoutExchange]: unable to publish a message " + err.Error())
	}

	log.Println("[PublishFanoutExchange]: publish message to queue")

	return nil
}

// PublishFanoutExchange publish message to queue who set routingKey match my exchange
func (p *producer) PublishTopicExchange(body []byte) error {
	err := p.publish(p.config.Exchange.Name, p.config.Exchange.RoutingKey, body)
	if err != nil {
		return errors.New("[PublishTopicExchange]: unable to publish a message " + err.Error())
	}

	log.Println("[PublishTopicExchange]: publish message to queue")

	return nil
}

// Close release the connection memory
func (p *producer) Close() {
	p.conn.Close()
	p.ch.Close()
}

func (p *producer) publish(exchange, routing string, body []byte) error {
	mandatory := false
	immediate := false

	bMessage := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/json",
		Body:         body,
	}

	return p.ch.Publish(exchange, routing, mandatory, immediate, bMessage)
}
