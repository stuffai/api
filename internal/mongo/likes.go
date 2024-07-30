package mongo

import (
	"context"

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
	ImageID *primitive.ObjectID
}

func initLikesCollection(ctx context.Context) error {
	_, err := likesCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{"image", 1}}, Options: options.Index().SetUnique(false)},
		{Keys: bson.D{{"user", 1}, {"image", 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return err
	}
	log.Info("mongo.initLikesCollection: success")
	return nil
}

// func InsertLike() error {
// 	//
// 	return nil
// }

// func DeleteLike() error {
// 	return nil
// }
