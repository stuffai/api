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
		{"userID", 1},
		{"state", 1},
		{"bucket", 1},
		{"dtModified", 1},
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

func findImages(ctx context.Context, query interface{}, opt *options.FindOptions) ([]*types.Image, error) {
	cur, err := imagesCollection().Find(ctx, query, opt)
	if err != nil {
		return nil, err
	}
	var imgs []*types.Image
	if err = cur.All(ctx, &imgs); err != nil {
		return nil, err
	}
	return imgs, nil
}

var orderDescending = bson.D{{"dtModified", -1}}

func FindImages(ctx context.Context) ([]*types.Image, error) {
	return findImages(ctx, bson.D{{"state", 1}}, options.Find().SetSort(orderDescending))
}

func FindImagesForUser(ctx context.Context, uid interface{}) ([]*types.Image, error) {
	return findImages(ctx, bson.D{{"userID", uid}}, options.Find().SetSort(orderDescending))
}
