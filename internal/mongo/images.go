package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stuff-ai/api/pkg/types"
)

func FindImagesAggregate(ctx context.Context, uid interface{}) ([]*types.Image, error) {
	log.Info("userID IN MONGO", uid)

	cur, err := jobsCollection().Aggregate(ctx, bson.A{
		bson.D{
			{"$match", bson.D{{"state", 1}}},
		},
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
			{"$lookup",
				bson.D{
					{"from", "craft_comments"},
					{"localField", "_id"},
					{"foreignField", "craftID"},
					{"as", "comments"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "craft_likes"},
					{"localField", "_id"},
					{"foreignField", "craftID"},
					{"as", "likes"},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "craft_likes"},
					{"let", bson.D{{"craftID", "$_id"}}},
					{"pipeline", mongo.Pipeline{
						{{"$match", bson.D{
							{"$expr", bson.D{
								{"$and", bson.A{
									bson.D{{"$eq", bson.A{"$userID", uid}}},
									bson.D{{"$eq", bson.A{"$craftID", "$$craftID"}}},
								}},
							}},
						}}},
					}},
					{"as", "user_like"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"user._id", "$userDocs._id"},
					{"user.ppBucket", "$userDocs.profile.ppBucket"},
					{"user.username", "$userDocs.username"},
					{"rank", 1},
					{"bucket", 1},
					{"dtModified", 1},
					{"title", "$promptDocs.title"},
					{"prompt", "$promptDocs.prompt"},
					{"_id", 1},
					{"nComments", bson.D{{"$size", "$comments"}}},
					{"nLikes", bson.D{{"$size", "$likes"}}},
					{"isLiked", bson.D{
						{"$eq", bson.A{
							bson.D{{"$size", "$user_like"}},
							1,
						}},
					}},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	images := []*types.Image{}
	if err := cur.All(ctx, &images); err != nil {
		return nil, err
	}
	return images, err
}

var viewProjection = bson.A{
	bson.D{
		{"$match", bson.D{{"state", 1}}},
	},
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
		{"$lookup",
			bson.D{
				{"from", "craft_comments"},
				{"localField", "_id"},
				{"foreignField", "craftID"},
				{"as", "comments"},
			},
		},
	},
	bson.D{
		{"$lookup",
			bson.D{
				{"from", "craft_likes"},
				{"localField", "_id"},
				{"foreignField", "craftID"},
				{"as", "likes"},
			},
		},
	},
	bson.D{
		{"$project",
			bson.D{
				{"user._id", "$userDocs._id"},
				{"user.ppBucket", "$userDocs.profile.ppBucket"},
				{"user.username", "$userDocs.username"},
				{"rank", 1},
				{"bucket", 1},
				{"dtModified", 1},
				{"title", "$promptDocs.title"},
				{"prompt", "$promptDocs.prompt"},
				{"_id", 1},
				{"nComments", bson.D{{"$size", "$comments"}}},
				{"nLikes", bson.D{{"$size", "$likes"}}},
			},
		},
	},
}

func createImagesView(ctx context.Context) error {
	log.Info("mongo.createImagesView")
	imagesCollection().Drop(ctx)
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
	return findImages(ctx, bson.D{}, options.Find().SetSort(orderDescending))
}

func FindImagesForUser(ctx context.Context, uid interface{}) ([]*types.Image, error) {
	return findImages(ctx, bson.D{{"user._id", uid}}, options.Find().SetSort(orderDescending))
}

func FindImageByID(ctx context.Context, craftID string) (*types.Image, error) {
	cid, _ := primitive.ObjectIDFromHex(craftID)
	craft := new(types.Image)
	return craft, imagesCollection().FindOne(ctx, bson.D{{"_id", cid}}).Decode(craft)
}
