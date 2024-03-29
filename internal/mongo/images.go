package mongo

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stuff-ai/api/pkg/types"
)

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

var orderDescending = bson.D{{"rank", -1}}

func FindImages(ctx context.Context) ([]*types.Image, error) {
	return findImages(ctx, bson.D{}, options.Find().SetSort(orderDescending))
}

func FindImagesForUser(ctx context.Context, uid interface{}) ([]*types.Image, error) {
	return findImages(ctx, bson.D{{"user._id", uid}}, options.Find().SetSort(orderDescending))
}

func FindImageRanksByIDs(ctx context.Context, ids [3]string) ([3]int, error) {
	out := [3]int{}
	idobj := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		obj, _ := primitive.ObjectIDFromHex(id)
		idobj[i] = obj
	}
	imgs, err := findImages(ctx, bson.D{{"_id", bson.D{{"$in", idobj}}}}, options.Find().SetProjection(bson.D{{"rank", 1}}))
	if err != nil {
		return out, err
	}
	// coerce to dict
	rankMap := map[string]int{}
	for _, img := range imgs {
		rankMap[img.ID] = img.Rank
	}
	// look up map and return ranks in order
	for i, id := range ids {
		out[i] = rankMap[id]
	}
	return out, nil
}

func UpdateImageRanks(ctx context.Context, rankMap map[string]int) error {
	updates := []mongo.WriteModel{}
	for id, rank := range rankMap {
		oid, _ := primitive.ObjectIDFromHex(id)
		updates = append(updates, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": oid}).SetUpdate(bson.M{"$set": bson.M{"rank": rank}}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := jobsCollection().BulkWrite(context.TODO(), updates, opts)
	if err != nil {
		return err
	}
	return nil

}

var getRankRandomSamplePipeline = mongo.Pipeline{
	{{"$group", bson.D{
		{"_id", "$user._id"},
		{"docs", bson.D{{"$push", "$$ROOT"}}},
	}}},
	{{"$unwind", "$docs"}},
	{{"$sample", bson.D{{"size", 3}}}},
	{{"$replaceRoot", bson.D{{"newRoot", "$docs"}}}},
}

func FindImagesForRank(ctx context.Context) ([]*types.Image, error) {
	cur, err := imagesCollection().Aggregate(ctx, getRankRandomSamplePipeline)
	if err != nil {
		log.WithError(err).Error("RandomPrompt: collection.Aggregate")
		return nil, err
	}
	var results []*types.Image
	if err = cur.All(ctx, &results); err != nil {
		log.WithError(err).Error("RandomPrompt: cursor.Decode")
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("RandomPrompt: no results")
	}
	return results, nil
}
