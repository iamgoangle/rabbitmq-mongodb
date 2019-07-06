package worker

import (
	"log"
	"time"

	"math/rand"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
)

// Worker behavior
type Worker interface {
	Processor(msgs <-chan amqp.Delivery)
}

type worker struct {
	db            *mongo.Database
	coll          *mongo.Collection
	qConn         *amqp.Connection
	delayProducer rabbitmq.Producer
}

const collection = "message_stat"

// New instance worker
func New(mgDB *mongo.Database, qConn *amqp.Connection) Worker {
	coll := mgDB.Collection(collection)

	delayQCfg := rabbitmq.ConfigProducer{
		Exchange: rabbitmq.ConfigExchange{
			Type:       rabbitmq.ExchangeTopic,
			Name:       "delay.exchange",
			RoutingKey: "delay.routing",
		},
		Queue: rabbitmq.ConfigQueue{
			Name: "delay.queue",
		},
	}
	delay, err := rabbitmq.NewProducer(qConn, delayQCfg)
	if err != nil {
		log.Panic(err)
	}

	return &worker{
		db:            mgDB,
		coll:          coll,
		delayProducer: delay,
	}
}

func (c *worker) Processor(msgs <-chan amqp.Delivery) {
	for m := range msgs {
		if random() == 0 {
			log.Println("Add to delay queue since service resource not done yet...")

			err := c.delayQueue(m.Body)
			if err != nil {
				// TODO: Add to retry exchange
				log.Println("Add to retry queue")
				continue
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

func (c *worker) delayQueue(msg []byte) error {
	err := c.delayProducer.Publish(msg)
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
