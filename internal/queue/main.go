package queue

import (
	"context"

	"cloud.google.com/go/pubsub"

	"github.com/stuff-ai/api/pkg/config"
)

var (
	topicIDGenerate string
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
}

func PublishGenerate(ctx context.Context, b []byte) error {
	msg := &pubsub.Message{Data: b}

	_, err := client.Topic(topicIDGenerate).Publish(ctx, msg).Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
