package queue

import (
	"context"

	"cloud.google.com/go/pubsub"

	"github.com/stuff-ai/api/pkg/config"
)

var (
	topicID string
	client  *pubsub.Client
)

func init() {
	ctx := context.Background()
	projectID := config.ProjectID()

	var err error
	if client, err = pubsub.NewClient(ctx, projectID); err != nil {
		panic(err)
	}

	topicID = config.PubSubTopicID()
}

func Publish(ctx context.Context, b []byte) error {
	msg := &pubsub.Message{Data: b}

	_, err := client.Topic(topicID).Publish(ctx, msg).Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
