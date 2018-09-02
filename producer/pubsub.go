package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		os.Exit(1)
	}
	return v
}

func newPubSubClient(ctx context.Context) (*pubsub.Client, error) {
	project := mustGetenv("GOOGLE_CLOUD_PROJECT")
	return pubsub.NewClient(ctx, project)
}

func getTopic(cli *pubsub.Client) *pubsub.Topic {
	topic := mustGetenv("PUBSUB_TOPIC")
	return cli.Topic(topic)
}

func createTopic(ctx context.Context, cli *pubsub.Client) (*pubsub.Topic, error) {
	topic := mustGetenv("PUBSUB_TOPIC")

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
		return fmt.Sprintf("https://%s.appspot.com/push", mustGetenv("GOOGLE_CLOUD_PROJECT"))
	}
}

func createSubs(ctx context.Context, cli *pubsub.Client, topic *pubsub.Topic, name string) (*pubsub.Subscription, error) {
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
