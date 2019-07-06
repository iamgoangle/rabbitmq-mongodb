package main

import (
	"context"
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/worker"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/mongodb"
	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
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
	w := worker.New(mgDBClient, rbConn)
	go w.Processor(msgs)

	<-forever
}

// func testInsert(coll *mongo.Collection) {
// 	insertData := &database.DeliveryMessage{
// 		RequestID:       "request:1234567",
// 		DeliveryCount:   100,
// 		UnDeliveryCount: 0,
// 		CreatedAt:       time.Now().Unix(),
// 	}

// 	insertResult, err := database.InsertMessage(coll, insertData)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	fmt.Println(insertResult)
// }

// func testUpdate(coll *mongo.Collection) {
// 	updateData := &database.DeliveryMessage{
// 		DeliveryCount: 50,
// 		UpdatedAt:     time.Now().Unix(),
// 	}

// 	filter := bson.D{
// 		{
// 			"requestId", "request:1234567",
// 		},
// 	}

// 	updateResult, err := database.UpdateMessage(coll, updateData, filter)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	fmt.Println(updateResult)
// }

// func testUpsert(coll *mongo.Collection) {
// 	// updateData := &database.DeliveryMessage{
// 	// 	DeliveryCount: 1000,
// 	// 	CreatedAt:     time.Now().Unix(),
// 	// }

// 	updateData := &database.DeliveryMessage{
// 		DeliveryCount: 500,
// 		UpdatedAt:     time.Now().Unix(),
// 	}

// 	filter := bson.D{
// 		{
// 			"requestId", "request:555555",
// 		},
// 	}

// 	updateResult, err := database.UpsertMessage(coll, updateData, filter)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	fmt.Println(updateResult)
// }

// func testFindOneAndUpsert(coll *mongo.Collection) {
// 	// updateData := &database.DeliveryMessage{
// 	// 	RequestID:       "request:555551",
// 	// 	DeliveryCount:   100,
// 	// 	UnDeliveryCount: 0,
// 	// 	CreatedAt:       time.Now().Unix(),
// 	// }

// 	updateData := &database.DeliveryMessage{
// 		DeliveryCount: 500,
// 		UpdatedAt:     time.Now().Unix(),
// 	}

// 	filter := bson.D{
// 		{
// 			"requestId", "request:555551",
// 		},
// 	}

// 	updateResult, err := database.FindOneAndUpsert(coll, updateData, filter)
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	r, err := updateResult.DecodeBytes()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	fmt.Println(string(r))
// }

// func testFindOneAndUpdate(coll *mongo.Collection) {
// 	// updateData := &database.DeliveryMessage{
// 	// 	RequestID:       "request:1000",
// 	// 	DeliveryCount:   100,
// 	// 	UnDeliveryCount: 0,
// 	// 	CreatedAt:       time.Now().Unix(),
// 	// }

// 	updateData := &database.DeliveryMessage{
// 		DeliveryCount: 1,
// 		UpdatedAt:     time.Now().Unix(),
// 	}

// 	filter := bson.D{
// 		{
// 			"requestId", "request:1000",
// 		},
// 	}

// 	updateResult, err := database.FindOneAndUpdate(coll, updateData, filter)
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	r, err := updateResult.DecodeBytes()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	fmt.Println(string(r))
// }
