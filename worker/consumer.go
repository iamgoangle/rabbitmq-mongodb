package worker

import (
	"log"
	"time"

	"math/rand"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Processor(msgs <-chan amqp.Delivery)
}

type consumer struct {
	client *mongo.Client
	coll   *mongo.Collection
	qConn  *amqp.Connection
	delay  rabbitmq.Producer
}

func NewConsumer(client *mongo.Client, db, collection string, qConn *amqp.Connection) Consumer {
	coll := client.Database(db).Collection(collection)

	delayConfig := rabbitmq.ConfigProducer{
		Exchange: rabbitmq.ConfigExchange{
			Type:       rabbitmq.ExchangeTopic,
			Name:       "retry.exchange",
			RoutingKey: "retry.routing",
		},
		Queue: rabbitmq.ConfigQueue{
			Name: "retry.queue",
		},
	}
	delay, err := rabbitmq.NewProducer(qConn, delayConfig)
	if err != nil {
		log.Panic(err)
	}

	return &consumer{
		client: client,
		coll:   coll,
		delay:  delay,
	}
}

func (c *consumer) Processor(msgs <-chan amqp.Delivery) {
	for m := range msgs {
		if random() == 0 {
			log.Println("Add to delay since service resouce not done yet...")

			err := c.delayQueue(m.Body)
			if err != nil {
				// TODO: Add to retry exchange
				log.Panic(err)
			}
			m.Reject(false)

			continue
		}

		log.Println("Received a message: %+v", string(m.Body))
		log.Println(m.Headers)
		log.Println("Done")
		m.Ack(false)
	}
}

func (c *consumer) delayQueue(msg []byte) error {
	err := c.delay.Publish(msg)
	if err != nil {
		return err
	}

	return nil
}

func random() int {
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 30

	mod := (rand.Intn(max-min) + min) % 2

	return mod
}
