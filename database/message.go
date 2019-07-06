package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message interface {
	Upsert(doc *DeliveryMessage, filter bson.D) error
	UpdateMessage(doc *DeliveryMessage, filter bson.D) (*mongo.UpdateResult, error)
	InsertMessage(doc interface{}) (*mongo.InsertOneResult, error)
}

type message struct {
	coll *mongo.Collection
}

const collection = "message_stat"

// New instance message db
func New(db *mongo.Database) Message {
	c := db.Collection(collection)

	return &message{
		coll: c,
	}
}

// UpdateMessage document
func (md *message) UpdateMessage(doc *DeliveryMessage, filter bson.D) (*mongo.UpdateResult, error) {
	bEncode, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}

	var bDecode bson.M
	err = bson.Unmarshal(bEncode, &bDecode)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": &bDecode,
		"$inc": bson.M{
			"callServiceCount": 1,
		},
	}

	opt := &options.UpdateOptions{}
	opt.SetUpsert(true)

	result, err := md.coll.UpdateOne(context.TODO(), filter, update, opt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Upsert method
func (md *message) Upsert(doc *DeliveryMessage, filter bson.D) error {
	bEncode, err := bson.Marshal(doc)
	if err != nil {
		return err
	}

	var bDecode bson.M
	err = bson.Unmarshal(bEncode, &bDecode)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": &bDecode,
		"$inc": bson.M{
			"callServiceCount": 1,
		},
	}

	opt := &options.FindOneAndUpdateOptions{}
	opt.SetUpsert(true)
	opt.SetReturnDocument(options.After)

	result := md.coll.FindOneAndUpdate(context.TODO(), filter, update, opt)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

// InsertMessage method
func (md *message) InsertMessage(doc interface{}) (*mongo.InsertOneResult, error) {
	result, err := md.coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}

	return result, nil
}
