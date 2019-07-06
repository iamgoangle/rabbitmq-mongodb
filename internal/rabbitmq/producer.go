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
	Publish(body []byte, headers amqp.Table) error

	// Release connection
	Close()
}

// ConfigProducer type
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

type execPublish func(body []byte, p *producer, headers amqp.Table) error

type producer struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config ConfigProducer

	execPublish execPublish
}

// NewProducer amqp producer actor
func NewProducer(conn *amqp.Connection, config ConfigProducer) (Producer, error) {
	var execPublish execPublish

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ExchangeDeclare(config.Exchange, ch)
	if err != nil {
		return nil, errors.New("[NewProducer] unable to declare exchange " + err.Error())
	}

	// queueExist := IsQueueExist(config.Queue.Name, ch)
	// if !queueExist {
	// 	_, err := QueueDeclare(config.Queue, ch)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	err = QueueBind(config.Queue, ch)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	switch config.Exchange.Type {
	case ExchangeDirect:
		execPublish = PublishDirectExchange
	case ExchangeFanout:
		execPublish = PublishFanoutExchange
	case ExchangeTopic:
		execPublish = PublishTopicExchange
	default:
		execPublish = PublishDefaultExchange
	}

	return &producer{
		ch:          ch,
		config:      config,
		conn:        conn,
		execPublish: execPublish,
	}, nil
}

// Publish factory method to exchange the message to queue by type
// TODO: To imrove init exchange type when NewProducer by pointer it to exechange type func
func (p *producer) Publish(body []byte, headers amqp.Table) error {
	return p.execPublish(body, p, headers)
}

// PublishDirectExchange handles default exchange
// Direct: A direct exchange delivers messages to queues based on a message routing key.
// In a direct exchange, the message is routed to the queues whose binding key exactly matches the routing key of the message.
func PublishDirectExchange(body []byte, p *producer, headers amqp.Table) error {
	err := p.publish(p.config.Exchange.Name, p.config.Exchange.RoutingKey, body)
	if err != nil {
		return errors.New("[PublishDirectExchange]: unable to publish a message " + err.Error())
	}

	log.Println("[PublishDirectExchange]: publish message to queue")

	return nil
}

// PublishDefaultExchange handles default exchange
func PublishDefaultExchange(body []byte, p *producer, headers amqp.Table) error {
	err := p.publish("", p.config.Queue.Name, body)
	if err != nil {
		return errors.New("[PublishDefaultExchange]: unable to publish a message " + err.Error())
	}

	log.Println("[PublishDefaultExchange]: publish message to queue")

	return nil
}

// PublishFanoutExchange broadcasts all the messages to queue that subscribe my exchange
func PublishFanoutExchange(body []byte, p *producer, headers amqp.Table) error {
	err := p.publish(p.config.Exchange.Name, "", body)
	if err != nil {
		return errors.New("[PublishFanoutExchange]: unable to publish a message " + err.Error())
	}

	log.Println("[PublishFanoutExchange]: publish message to queue")

	return nil
}

// PublishTopicExchange publish message to queue who set routingKey match my exchange
func PublishTopicExchange(body []byte, p *producer, headers amqp.Table) error {
	err := p.publish(p.config.Exchange.Name, p.config.Exchange.RoutingKey, body, headers)
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

func (p *producer) publish(exchange, routing string, body []byte, headers ...amqp.Table) error {
	mandatory := false
	immediate := false

	var h amqp.Table
	if len(headers) > 0 {
		h = headers[0]
	}

	bMessage := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/json",
		Body:         body,
		Headers:      h,
	}

	return p.ch.Publish(exchange, routing, mandatory, immediate, bMessage)
}
