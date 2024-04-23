package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stuff-ai/api/pkg/types"
)

func commentsCollection() *mongo.Collection {
	return db().Collection("craft_comments")
}

func FindComments(ctx context.Context, cid string) ([]*types.Comment, error) {
	coid, _ := primitive.ObjectIDFromHex(cid)
	cur, err := commentsCollection().Aggregate(ctx, bson.A{
		bson.D{
			{"$match", bson.D{{"craftID", coid}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "userID"},
				{"foreignField", "_id"},
				{"as", "user"},
			}},
		},
		bson.D{
			{"$unwind", "$user"},
		},
		bson.D{
			{"$project", bson.D{
				{"text", 1},
				{"userID", 1},
				{"username", "$user.username"},
				{"ppBucket", "$user.profile.ppBucket"},
				{"dtCreated", 1},
			}},
		},
	})
	if err != nil {
		return nil, err
	}

	comments := []*types.Comment{}
	if err := cur.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, err
}

func InsertComment(ctx context.Context, uid interface{}, cid, text string) error {
	coid, _ := primitive.ObjectIDFromHex(cid)
	_, err := commentsCollection().InsertOne(ctx, bson.D{
		{"craftID", coid}, {"userID", uid}, {"text", text}, {"dtCreated", time.Now()},
	})
	return err
}
