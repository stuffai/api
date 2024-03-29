package mongo

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stuff-ai/api/pkg/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ranksCollection() *mongo.Collection {
	return db().Collection("ranks")
}

type rankedImage struct {
	ID    primitive.ObjectID `json:"id" bson:"id"`
	Delta int                `json:"delta" bson:"delta"`
	Prev  int                `json:"prev" bson:"prev"`
}

type rank struct {
	User      primitive.ObjectID `json:"user" bson:"user"`
	First     rankedImage        `json:"first" bson:"first"`
	Second    rankedImage        `json:"second" bson:"second"`
	Third     rankedImage        `json:"third" bson:"third"`
	DTCreated time.Time          `json:"dtCreated" bson:"dtCreated"`
}

func InsertRank(ctx context.Context, uid interface{}, ranks [3]string, deltas, scores [3]int) error {
	first, _ := primitive.ObjectIDFromHex(ranks[0])
	second, _ := primitive.ObjectIDFromHex(ranks[1])
	third, _ := primitive.ObjectIDFromHex(ranks[2])
	r := &rank{
		User:      uid.(primitive.ObjectID),
		First:     rankedImage{ID: first, Delta: deltas[0], Prev: scores[0]},
		Second:    rankedImage{ID: second, Delta: deltas[1], Prev: scores[1]},
		Third:     rankedImage{ID: third, Delta: deltas[2], Prev: scores[2]},
		DTCreated: time.Now(),
	}
	_, err := ranksCollection().InsertOne(ctx, r)
	return err
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
		updates = append(updates, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": oid}).SetUpdate(bson.M{"$inc": bson.M{"rank": rank}}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := jobsCollection().BulkWrite(context.TODO(), updates, opts)
	if err != nil {
		return err
	}
	return nil

}
