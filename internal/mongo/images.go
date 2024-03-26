package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stuff-ai/api/pkg/types"
)

var viewProjection = bson.A{
	bson.D{
		{"$lookup",
			bson.D{
				{"from", "prompts"},
				{"localField", "promptID"},
				{"foreignField", "_id"},
				{"as", "promptDocs"},
			},
		},
	},
	bson.D{{"$unwind", "$promptDocs"}},
	bson.D{
		{"$lookup",
			bson.D{
				{"from", "users"},
				{"localField", "userID"},
				{"foreignField", "_id"},
				{"as", "userDocs"},
			},
		},
	},
	bson.D{{"$unwind", "$userDocs"}},
	bson.D{
		{"$project",
			bson.D{
				{"user._id", "$userDocs._id"},
				{"user.ppURL", "$userDocs.profile.ppURL"},
				{"user.username", "$userDocs.username"},
				{"state", 1},
				{"bucket", 1},
				{"dtModified", 1},
				{"title", "$promptDocs.title"},
				{"prompt", "$promptDocs.prompt"},
			},
		},
	},
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
	if imgs == nil {
		imgs = []*types.Image{}
	}
	return imgs, nil
}

var orderDescending = bson.D{{"dtModified", -1}}

func FindImages(ctx context.Context) ([]*types.Image, error) {
	return findImages(ctx, bson.D{{"state", 1}}, options.Find().SetSort(orderDescending))
}

func FindImagesForUser(ctx context.Context, uid interface{}) ([]*types.Image, error) {
	return findImages(ctx, bson.D{{"state", 1}, {"user._id", uid}}, options.Find().SetSort(orderDescending))
}
