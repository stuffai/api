package queue

import (
	"context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	projectID = "stuffai-dev"
	topicID   = "stuffai-dev"
	client    *pubsub.Client
)

func init() {
	ctx := context.Background()

	var err error
	if client, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("/home/b/.gcloud/stuffai-dev-52b1d6089a1d.json")); err != nil {
		panic(err)
	}
}

func Publish(ctx context.Context, b []byte) error {
	msg := &pubsub.Message{Data: b}

	_, err := client.Topic(topicID).Publish(ctx, msg).Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
