package mongo

import (
	"context"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stuff-ai/api/pkg/config"
)

var client *mongo.Client

func init() {
	uri := config.MongoURI()

	ctx := context.Background()
	var err error
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri)); err != nil {
		panic("mongo: failed to init: " + err.Error())
	}

	// check if views are initialized
	collections, err := db().ListCollectionNames(ctx, bson.D{}, nil)
	if err != nil {
		panic("mongo: failed to list collections: " + err.Error())
	}

	if !slices.Contains(collections, "images") {
		if err = createImagesView(ctx); err != nil {
			panic("mongo: failed to initialize images view: " + err.Error())
		}
	}

	// ensure indicies
	if err := initFriendsCollection(ctx); err != nil {
		panic("mongo: failed to ensure friends collection index: " + err.Error())
	}
	if err := initNotificationsCollection(ctx); err != nil {
		panic("mongo: failed to ensure notifs collection index: " + err.Error())
	}
}

func Shutdown() {
	client.Disconnect(context.Background())
}

func db() *mongo.Database {
	return client.Database("stuffai")
}
