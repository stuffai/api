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

func notificationsCollection() *mongo.Collection {
	return db().Collection("notifications")
}

func initNotificationsCollection(ctx context.Context) error {
	_, err := notificationsCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"userID", 1}}, Options: options.Index().SetUnique(false),
	})
	if err != nil {
		return err
	}
	log.Info("mongo.initNotificationsCollection: success")
	return nil
}

func InsertNotification(ctx context.Context, kind types.NotificationKind, data, uid interface{}) (primitive.ObjectID, error) {
	uidS, isString := uid.(string)
	if isString {
		uid, _ = primitive.ObjectIDFromHex(uidS)
	}
	result, err := notificationsCollection().InsertOne(ctx, bson.D{
		{"userID", uid}, {"kind", kind}, {"data", data}, {"read", false}, {"dtCreated", time.Now()},
	})
	return result.InsertedID.(primitive.ObjectID), err
}

func GetNotifications(ctx context.Context, uid interface{}) ([]*types.Notification, error) {
	cur, err := notificationsCollection().Find(ctx, bson.D{{"userID", uid}}, options.Find().SetSort(bson.D{{"dtCreated", -1}}))
	if err != nil {
		return nil, err
	}

	notifs := []*types.Notification{}
	if err := cur.All(ctx, &notifs); err != nil {
		return nil, err
	}

	return notifs, nil
}

func DeleteFriendRequestNotification(ctx context.Context, uid, fuid interface{}) error {
	_, err := notificationsCollection().DeleteOne(ctx, bson.D{{"userID", uid}, {"data.id", fuid}})
	return err
}

func UpdateNotificationRead(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := notificationsCollection().UpdateByID(ctx, oid, bson.D{{"$set", bson.D{{"read", true}}}})
	return err
}
