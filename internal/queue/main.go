package queue

import (
	"context"

	"cloud.google.com/go/pubsub"

	"github.com/stuff-ai/api/pkg/config"
)

var (
	topicIDGenerate string
	topicIDNotify   string
	client          *pubsub.Client
)

func init() {
	ctx := context.Background()
	projectID := config.ProjectID()

	var err error
	if client, err = pubsub.NewClient(ctx, projectID); err != nil {
		panic(err)
	}

	topicIDGenerate = config.PubSubTopicIDGenerate()
	topicIDNotify = config.PubSubTopicIDNotify()
}

func PublishGenerate(ctx context.Context, b []byte) error {
	msg := &pubsub.Message{Data: b}

	_, err := client.Topic(topicIDGenerate).Publish(ctx, msg).Get(ctx)
	if err != nil {
		return err
	}
	return nil
}

func PublishNotify(ctx context.Context, b []byte) error {
	msg := &pubsub.Message{Data: b}

	_, err := client.Topic(topicIDNotify).Publish(ctx, msg).Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
