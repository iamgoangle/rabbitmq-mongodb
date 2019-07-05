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
	testUpdateOne(coll)

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func testInsertOne(coll *mongo.Collection) {
	insertData := &database.DeliveryMessage{
		RequestID:       "request:1234567",
		DeliveryCount:   100,
		UnDeliveryCount: 0,
		CreatedAt:       time.Now().Unix(),
	}

	insertResult, err := database.InsertOne(coll, insertData)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(insertResult)
}

func testUpdateOne(coll *mongo.Collection) {
	updateData := &database.DeliveryMessage{
		DeliveryCount: 250,
		UpdatedAt:     time.Now().Unix(),
	}

	filter := bson.D{
		{
			"requestId", "request:1234567",
		},
	}

	updateResult, err := database.UpdateOne(coll, updateData, filter)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(updateResult)
}
