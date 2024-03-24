package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func jobsCollection() *mongo.Collection {
	return db().Collection("jobs")
}

func InsertJob(ctx context.Context, promptID string) (string, error) {
	promptOID, err := primitive.ObjectIDFromHex(promptID)
	if err != nil {
		return "", err
	}

	result, err := jobsCollection().InsertOne(
		ctx,
		bson.D{{"promptID", promptOID}, {"state", 0}},
	)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
