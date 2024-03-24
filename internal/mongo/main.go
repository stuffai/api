package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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

func promptsCollection() *mongo.Collection {
	return db().Collection("prompts")
}

func AddPrompt(ctx context.Context, title, prompt string) error {
	_, err := promptsCollection().InsertOne(
		ctx,
		bson.D{{"title", title}, {"prompt", prompt}},
	)
	return err
}
