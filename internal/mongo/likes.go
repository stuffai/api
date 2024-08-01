package mongo

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func likesCollection() *mongo.Collection {
	return db().Collection("craft_likes")
}

type like struct {
	ID      *primitive.ObjectID `bson:"_id"`
	UserID  *primitive.ObjectID
	CraftID *primitive.ObjectID
}

func initLikesCollection(ctx context.Context) error {
	_, err := likesCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"craftID", 1}}, Options: options.Index().SetUnique(false)},
		{Keys: bson.D{{"userID", 1}, {"craftID", 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return err
	}
	log.Info("mongo.initLikesCollection: success")
	return nil
}

func bsonDUserLikeCraft(coid primitive.ObjectID, uid interface{}) bson.D {
	return bson.D{
		{"craftID", coid}, {"userID", uid}, {"dtCreated", time.Now()},
	}
}

func bsonDDeleteUserLikeCraft(coid primitive.ObjectID, uid interface{}) bson.D {
	return bson.D{
		{"craftID", coid},
		{"userID", uid},
	}
}

func InsertLike(ctx context.Context, uid interface{}, cid string) error {
	coid, _ := primitive.ObjectIDFromHex(cid)
	_, err := likesCollection().InsertOne(ctx, bsonDUserLikeCraft(coid, uid))
	return err
}

func DeleteLike(ctx context.Context, uid interface{}, cid string) error {
	coid, _ := primitive.ObjectIDFromHex(cid)
	_, err := likesCollection().DeleteOne(ctx, bsonDDeleteUserLikeCraft(coid, uid))
	return err
}
