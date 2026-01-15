package scanning

import (
	"context"

	"cloud.google.com/go/pubsub"
)

func CreateTopic(ctx context.Context, client *pubsub.Client, id string) (topic *pubsub.Topic, err error) {
	topic = client.Topic(id)
	exists, err := topic.Exists(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		topic, err = client.CreateTopic(ctx, id)
		if err != nil {
			panic(err)
		}
	}

	return
}

