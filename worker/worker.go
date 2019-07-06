package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/iamgoangle/rabbitmq-mongodb/database"
	"go.mongodb.org/mongo-driver/bson"

	"math/rand"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/streadway/amqp"
)

// Worker behavior
type Worker interface {
	Processor(msgs <-chan amqp.Delivery)
}

type worker struct {
	msgDB         database.Message
	qConn         *amqp.Connection
	delayProducer rabbitmq.Producer
}

const collection = "message_stat"

// New instance worker
func New(msgDB database.Message, qConn *amqp.Connection) Worker {
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
		msgDB:         msgDB,
		delayProducer: delay,
	}
}

// Processor consume job to execute service resource, if deliveryCount is equal 0
// update data to mongo, if so trying to publish delay queue and try it again
func (c *worker) Processor(msgs <-chan amqp.Delivery) {
	for m := range msgs {
		undeliveryCount := deliveryCountService()

		log.Println("undeliveryCount ", undeliveryCount)

		var dataInQueue database.DeliveryMessage
		err := json.Unmarshal(m.Body, &dataInQueue)
		if err != nil {
			continue
		}
		dataInQueue.UnDeliveryCount = undeliveryCount

		if !dataInQueue.IsDone() {
			log.Println("Add to delay queue since service resource not done yet...")

			err := c.addJobToDelayQueue(m.Body)
			if err != nil {
				log.Println("Add to retry queue...")
				continue
			}
			m.Reject(false)
			c.updateDB(dataInQueue, m)
			continue
		}

		log.Println("Done")
		m.Ack(false)

		c.updateDB(dataInQueue, m)
	}
}

func (c *worker) addJobToDelayQueue(msg []byte) error {
	headers := amqp.Table{
		"updatedAt": time.Now().Unix(),
	}
	err := c.delayProducer.Publish(msg, headers)
	if err != nil {
		return err
	}

	return nil
}

func (c *worker) updateDB(d database.DeliveryMessage, m amqp.Delivery) error {
	var createdAt int64
	createdAt = time.Now().Unix()

	if len(m.Headers) != 0 {
		d.UpdatedAt = m.Headers["updatedAt"].(int64)
	}

	if d.UpdatedAt == 0 {
		d.CreatedAt = createdAt
		d.UpdatedAt = createdAt
	}

	fmt.Println("database.DeliveryMessage")
	fmt.Printf("%+v", d)

	filter := bson.D{
		{
			"requestId", d.RequestID,
		},
	}

	_, err := c.msgDB.UpdateMessage(&d, filter)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("update mongodb")

	return nil
}

func deliveryCountService() int {
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 30

	mod := (rand.Intn(max-min) + min) % 2

	return mod
}
