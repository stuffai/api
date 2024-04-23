package queue

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
	log "github.com/sirupsen/logrus"

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

func PublishNotifyMany(ctx context.Context, b [][]byte) {
	n := len(b)
	t := client.Topic(topicIDNotify)
	t.PublishSettings.CountThreshold = n
	t.PublishSettings.ByteThreshold = n * 13
	results := make([]*pubsub.PublishResult, n)
	for i, msg := range b {
		results[i] = t.Publish(ctx, &pubsub.Message{Data: msg})
	}
	wg := new(sync.WaitGroup)
	for _, res := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := res.Get(ctx)
			if err != nil {
				log.WithError(err).Error("queue.PublishNotifyMany")
			}
		}()
	}
	wg.Wait()
}
