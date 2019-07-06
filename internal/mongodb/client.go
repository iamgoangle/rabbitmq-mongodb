package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database() *mongo.Database
}

type Config struct {
	Host string
	Db   string
}

type mongoClient struct {
	client *mongo.Client
	config Config
}

// NewMongoClient initial mongo connection
func NewMongoClient(c Config) MongoClient {
	clientOptions := options.Client().ApplyURI(c.Host)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// return client.Database(c.Db)
	return &mongoClient{
		client: client,
		config: c,
	}
}

func (c *mongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.client.Ping(ctx, nil)
}

func (c *mongoClient) Database() *mongo.Database {
	return c.client.Database(c.config.Db)
}
