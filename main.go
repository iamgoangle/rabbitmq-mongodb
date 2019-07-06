package main

import (
	"context"
	"log"
	"time"

	"github.com/iamgoangle/rabbitmq-mongodb/database"
	"github.com/iamgoangle/rabbitmq-mongodb/internal/mongodb"
	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
	"github.com/iamgoangle/rabbitmq-mongodb/worker"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	mongoHost = "mongodb://localhost:27017"
	db        = "golf-message"

	rabbitMqHost = "amqp://admin:1234@localhost:5672"
	rabbitMqType = "standalone"
)

func main() {
	// MongoDB
	mgConfig := mongodb.Config{
		Host: mongoHost,
		Db:   db,
	}
	mgClient := mongodb.NewMongoClient(mgConfig)
	err := mgClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic(err)
	}
	mgDBClient := mgClient.Database()
	msgDB := database.New(mgDBClient)

	// RabbitMQ
	rbConn, err := rabbitmq.NewConnection(rabbitmq.ConfigConnection{
		Type: rabbitMqType,
		Url:  rabbitMqHost,
	})
	if err != nil {
		log.Fatalln("[main]: unable to connect RabbitMQ %+v", err)
	}

	worker.InitExchange(rbConn)
	worker.InitQueue(rbConn)

	workQCfg := rabbitmq.ConfigConsumer{
		Exchange: rabbitmq.ConfigExchange{
			Type: rabbitmq.ExchangeTopic,
			Name: "work.exchange",
		},
		Queue: rabbitmq.ConfigQueue{
			Name: "work.queue",
			Bind: rabbitmq.ConfigQueueBind{
				ExchangeName: "work.exchange",
				RoutingKey:   "work.routing.#",
			},
		},
	}

	consumer, err := rabbitmq.NewConsumer(rbConn, workQCfg)
	if err != nil {
		log.Panic(err)
	}

	msgs, err := consumer.WorkerProcessor()
	if err != nil {
		log.Panic(err)
	}

	forever := make(chan bool)

	// worker
	w := worker.New(msgDB, rbConn)
	go w.Processor(msgs)

	<-forever
}

func testFindOneAndUpsert(msgDB database.Message) {
	// updateData := &database.DeliveryMessage{
	// 	RequestID:       "1552727174:123456789:chunk:000001",
	// 	ChannelID:       "1549045957",
	// 	LineRequestID:   "2e1d92b8-9e11-4d39-8c1c-2ef0f40934c7",
	// 	DeliveryCount:   10,
	// 	UnDeliveryCount: 100,
	// 	CreatedAt:       time.Now().Unix(),
	// }

	updateData := &database.DeliveryMessage{
		DeliveryCount:   0,
		UnDeliveryCount: 0,
		UpdatedAt:       time.Now().Unix(),
	}

	filter := bson.D{
		{
			"requestId", "1552727174:123456789:chunk:000001",
		},
	}

	_, err := msgDB.UpdateMessage(updateData, filter)
	if err != nil {
		log.Panicln(err)
	}
}
