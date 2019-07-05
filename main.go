package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/iamgoangle/rabbitmq-mongodb/database"
)

const (
	mongoHost  = "mongodb://localhost:27017"
	db         = "golf-message"
	collection = "message_stat"
)

func main() {
	clientOptions := options.Client().ApplyURI(mongoHost)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	coll := client.Database(db).Collection(collection)
	// testInsert(coll)
	// testUpdate(coll)
	// testUpsert(coll)
	testFindOneAndUpsert(coll)

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func testInsert(coll *mongo.Collection) {
	insertData := &database.DeliveryMessage{
		RequestID:       "request:1234567",
		DeliveryCount:   100,
		UnDeliveryCount: 0,
		CreatedAt:       time.Now().Unix(),
	}

	insertResult, err := database.InsertMessage(coll, insertData)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(insertResult)
}

func testUpdate(coll *mongo.Collection) {
	updateData := &database.DeliveryMessage{
		DeliveryCount: 50,
		UpdatedAt:     time.Now().Unix(),
	}

	filter := bson.D{
		{
			"requestId", "request:1234567",
		},
	}

	updateResult, err := database.UpdateMessage(coll, updateData, filter)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(updateResult)
}

func testUpsert(coll *mongo.Collection) {
	// updateData := &database.DeliveryMessage{
	// 	DeliveryCount: 1000,
	// 	CreatedAt:     time.Now().Unix(),
	// }

	updateData := &database.DeliveryMessage{
		DeliveryCount: 500,
		UpdatedAt:     time.Now().Unix(),
	}

	filter := bson.D{
		{
			"requestId", "request:555555",
		},
	}

	updateResult, err := database.UpsertMessage(coll, updateData, filter)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(updateResult)
}

func testFindOneAndUpsert(coll *mongo.Collection) {
	// updateData := &database.DeliveryMessage{
	// 	RequestID:       "request:555551",
	// 	DeliveryCount:   100,
	// 	UnDeliveryCount: 0,
	// 	CreatedAt:       time.Now().Unix(),
	// }

	updateData := &database.DeliveryMessage{
		DeliveryCount: 500,
		UpdatedAt:     time.Now().Unix(),
	}

	filter := bson.D{
		{
			"requestId", "request:555551",
		},
	}

	updateResult, err := database.FindOneAndUpsert(coll, updateData, filter)
	if err != nil {
		log.Panicln(err)
	}

	r, err := updateResult.DecodeBytes()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(string(r))
}
