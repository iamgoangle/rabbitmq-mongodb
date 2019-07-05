package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertMessage document
func InsertMessage(coll *mongo.Collection, doc interface{}) (*mongo.InsertOneResult, error) {
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateMessage document
func UpdateMessage(coll *mongo.Collection, doc *DeliveryMessage, filter bson.D) (*mongo.UpdateResult, error) {
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
	opt.SetUpsert(false)

	result, err := coll.UpdateOne(context.TODO(), filter, update, opt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpsertMessage method
func UpsertMessage(coll *mongo.Collection, doc *DeliveryMessage, filter bson.D) (*mongo.UpdateResult, error) {
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

	result, err := coll.UpdateOne(context.TODO(), filter, update, opt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindOneAndUpsert method
func FindOneAndUpsert(coll *mongo.Collection, doc *DeliveryMessage, filter bson.D) (*mongo.SingleResult, error) {
	bEncode, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}

	opt := &options.FindOneAndReplaceOptions{}
	opt.SetUpsert(true)
	opt.SetReturnDocument(options.After)

	result := coll.FindOneAndReplace(context.TODO(), filter, bEncode, opt)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}

// FindOneAndUpdate method
func FindOneAndUpdate(coll *mongo.Collection, doc *DeliveryMessage, filter bson.D) (*mongo.SingleResult, error) {
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

	opt := &options.FindOneAndUpdateOptions{}
	opt.SetUpsert(true)
	opt.SetReturnDocument(options.After)

	result := coll.FindOneAndUpdate(context.TODO(), filter, update, opt)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}
