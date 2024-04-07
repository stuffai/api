package mongo

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stuff-ai/api/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func friendsCollection() *mongo.Collection {
	return db().Collection("users_friends")
}

type friend struct {
	From *primitive.ObjectID
	To   *primitive.ObjectID
}

func initFriendsCollection(ctx context.Context) error {
	_, err := friendsCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"from", 1}}, Options: options.Index().SetUnique(false)},
		{Keys: bson.D{{"from", 1}, {"to", 1}}, Options: options.Index().SetUnique(true)},
	})
	return err
}

func bsonDFriendsFromTo(fromUserID, toUserID interface{}) bson.D {
	return bson.D{{"from", fromUserID}, {"to", toUserID}}
}

func FindFriends(ctx context.Context, uid interface{}) ([]*types.UserProfile, error) {
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"from", uid}, {"dtAccepted", bson.D{{"$ne", nil}}}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "to"},
			{"foreignField", "_id"},
			{"as", "friend"}}}},
		bson.D{{"$unwind", "$friend"}},
		bson.D{{"$project", bson.D{
			{"username", "$friend.username"},
			{"ppBucket", "$friend.profile.ppBucket"},
		}}},
	}
	friends := []*types.UserProfile{}
	cur, err := friendsCollection().Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &friends); err != nil {
		return nil, err
	}
	return friends, nil
}

func FindFriendRequests(ctx context.Context, uid interface{}) ([]*types.UserProfile, error) {
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"to", uid}, {"dtAccepted", nil}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "from"},
			{"foreignField", "_id"},
			{"as", "friend"}}}},
		bson.D{{"$unwind", "$friend"}},
		bson.D{{"$project", bson.D{
			{"username", "$friend.username"},
			{"ppBucket", "$friend.profile.ppBucket"},
		}}},
	}
	friends := []*types.UserProfile{}
	cur, err := friendsCollection().Aggregate(ctx, pipeline)
	if err != nil {
		log.WithError(err).Error("aggregate")
		return nil, err
	}
	if err := cur.All(ctx, &friends); err != nil {
		log.WithError(err).Error("cur.all")
		return nil, err
	}
	return friends, nil
}

func InsertFriendRequest(ctx context.Context, fromUID, toUID interface{}) error {
	_, err := friendsCollection().InsertOne(ctx,
		bson.D{{"from", fromUID}, {"to", toUID}, {"dtCreated", time.Now()}, {"dtAccepted", nil}},
	)
	return err
}

func ExistsFriendRequest(ctx context.Context, fromUserID, toUserID interface{}) (bool, error) {
	obj := bson.M{}
	err := friendsCollection().FindOne(ctx, bsonDFriendsFromTo(fromUserID, toUserID), options.FindOne()).Decode(&obj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func AcceptFriendRequest(ctx context.Context, fromUID, toUID interface{}) error {
	models := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(bsonDFriendsFromTo(toUID, fromUID)).SetUpdate(bson.D{{"$set", bson.D{{"dtAccepted", time.Now()}}}}),
		mongo.NewInsertOneModel().SetDocument(bson.D{{"from", fromUID}, {"to", toUID}, {"dtCreated", time.Now()}, {"dtAccepted", time.Now()}}),
	}
	_, err := friendsCollection().BulkWrite(ctx, models, options.BulkWrite().SetOrdered(true))
	return err
}

func RejectFriendRequest(ctx context.Context, fromUserID interface{}, toUserID interface{}) error {
	_, err := friendsCollection().DeleteOne(ctx, bsonDFriendsFromTo(fromUserID, toUserID))
	return err
}
