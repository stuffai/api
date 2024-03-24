package mongo

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stuff-ai/api/pkg/types"
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

func AddPrompt(ctx context.Context, prompt *types.Prompt) error {
	_, err := promptsCollection().InsertOne(
		ctx,
		bson.D{{"title", prompt.Title}, {"prompt", prompt.Prompt}},
	)
	return err
}

func RandomPrompt(ctx context.Context) (*types.Prompt, error) {
	filter := mongo.Pipeline{bson.D{{"$sample", bson.D{{"size", 1}}}}}
	cur, err := promptsCollection().Aggregate(ctx, filter)
	if err != nil {
		log.WithError(err).Error("RandomPrompt: collection.Aggregate")
		return nil, err
	}
	var results []*types.Prompt
	if err = cur.All(ctx, &results); err != nil {
		log.WithError(err).Error("RandomPrompt: cursor.Decode")
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("RandomPrompt: no results")
	}
	return results[0], nil
}
