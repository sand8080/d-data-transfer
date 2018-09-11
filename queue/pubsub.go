package queue

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sand8080/d-data-transfer/env"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewPubSubClient(ctx context.Context) (*pubsub.Client, error) {
	project := env.MustGetenv("GOOGLE_CLOUD_PROJECT")
	return pubsub.NewClient(ctx, project)
}

func GetTopic(cli *pubsub.Client) *pubsub.Topic {
	topic := env.MustGetenv("PUBSUB_TOPIC")
	return cli.Topic(topic)
}

func CreateTopic(ctx context.Context, cli *pubsub.Client) (*pubsub.Topic, error) {
	topic := env.MustGetenv("PUBSUB_TOPIC")

	// Create a topic to subscribe to.
	t := cli.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return t, nil
	}

	return cli.CreateTopic(ctx, topic)
}

func getPushEndpoint(ctx context.Context) string {
	log.Debugf(ctx, "Checking if in dev environment")
	if appengine.IsDevAppServer() {
		return "http://localhost:8080/push"
	} else {
		return fmt.Sprintf("https://%s.appspot.com/push", env.MustGetenv("GOOGLE_CLOUD_PROJECT"))
	}
}

func CreateSubs(ctx context.Context, cli *pubsub.Client, topic *pubsub.Topic, name string) (*pubsub.Subscription, error) {
	pushCfg := pubsub.PushConfig{Endpoint: getPushEndpoint(ctx)}
	sub, err := cli.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:       topic,
		PushConfig:  pushCfg,
		AckDeadline: 20 * time.Second,
	})
	s, ok := status.FromError(err)
	if !ok || s.Code() != codes.AlreadyExists {
		if err != nil {
			log.Errorf(ctx, "Subscrbtion creation failed: %v", err)
			return nil, err
		}
	}
	return sub, nil
}
