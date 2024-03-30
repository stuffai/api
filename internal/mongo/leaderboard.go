package mongo

import (
	"context"

	"github.com/stuff-ai/api/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var leaderboardProjection = bson.A{
	bson.D{{"$lookup",
		bson.D{
			{"from", "images"},
			{"localField", "_id"},
			{"foreignField", "user._id"},
			{"as", "userImages"}}},
	},
	bson.D{{"$unwind", "$userImages"}},
	bson.D{{"$sort",
		bson.D{
			{"userImages.rank", -1},
		}}},
	bson.D{{"$group",
		bson.D{
			{"_id", "$_id"},
			{"userData", bson.D{{"$first", "$$ROOT"}}},
			{"topImage", bson.D{{"$first", "$userImages"}}}}}},
	bson.D{{"$project",
		bson.D{
			{"_id", "$userData._id"},
			{"username", "$userData.username"},
			{"rank", "$topImage.rank"},
			{"ppBucket", "$userData.profile.bucket"}}}},
	bson.D{{"$limit", 10}},
}

func FindLeaderboard(ctx context.Context) ([]*types.LeaderboardEntry, error) {
	cur, err := usersCollection().Aggregate(ctx, leaderboardProjection, options.Aggregate())
	if err != nil {
		return nil, err
	}
	entries := []*types.LeaderboardEntry{}
	if err := cur.All(ctx, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
