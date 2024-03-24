package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	uri := "mongodb://192.168.63.29:27017/stuffai"
	var err error
	if client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri)); err != nil {
		panic("mongo: failed to init: " + err.Error())
	}
}

func Shutdown() {
	client.Disconnect(context.Background())
}

func db() *mongo.Database {
	return client.Database("stuffai")
}
