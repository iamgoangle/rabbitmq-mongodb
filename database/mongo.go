package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne document
func InsertOne(coll *mongo.Collection, doc interface{}) (*mongo.InsertOneResult, error) {
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateOne document
func UpdateOne(coll *mongo.Collection, doc interface{}, filter bson.D) (*mongo.UpdateResult, error) {
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
	}

	opt := &options.UpdateOptions{}
	opt.SetUpsert(false)

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
