package mongo

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stuff-ai/api/pkg/types"
)

func promptsCollection() *mongo.Collection {
	return db().Collection("prompts")
}

func AddPrompt(ctx context.Context, prompt *types.Prompt) (string, error) {
	result, err := promptsCollection().InsertOne(
		ctx,
		bson.D{{"title", prompt.Title}, {"prompt", prompt.Prompt}},
	)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
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
