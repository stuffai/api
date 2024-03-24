package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stuff-ai/api/pkg/types"
)

var viewProjection = mongo.Pipeline{
	bson.D{{"$lookup", bson.D{
		{"from", "prompts"},
		{"localField", "promptID"},
		{"foreignField", "_id"},
		{"as", "promptDocs"},
	}}},
	bson.D{{"$project", bson.D{
		{"bucket", 1},
		{"dtCreated", 1},
		{"title", "$promptDocs.title"},
		{"prompt", "$promptDocs.prompt"},
	}}},
	bson.D{{"$unwind", "$title"}},
	bson.D{{"$unwind", "$prompt"}},
}

func createImagesView(ctx context.Context) error {
	log.Info("mongo.createImagesView")
	return db().CreateView(ctx, "images", "jobs", viewProjection, nil)
}

func imagesCollection() *mongo.Collection {
	return db().Collection("images")
}

func FindImages(ctx context.Context) ([]*types.Image, error) {
	cur, err := imagesCollection().Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{"dtCreated", -1}}))
	if err != nil {
		return nil, err
	}
	var imgs []*types.Image
	if err = cur.All(ctx, &imgs); err != nil {
		return nil, err
	}
	return imgs, nil
}
