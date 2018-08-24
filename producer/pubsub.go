package main

import (
	"golang.org/x/net/context"
	syslog "log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		syslog.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func newPubSubClient() *pubsub.Client {
	project := mustGetenv("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		syslog.Fatalf("Could not create pubsub Client: %v", err)
	}

	return client
}

func createTopic(c *pubsub.Client) *pubsub.Topic {
	ctx := context.Background()
	topic := mustGetenv("PUBSUB_TOPIC")

	// Create a topic to subscribe to.
	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		syslog.Fatal(err)
	}
	if ok {
		return t
	}

	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		syslog.Fatalf("Failed to create the topic: %v", err)
	}
	return t
}

func createSubs(client *pubsub.Client, subName string, topic *pubsub.Topic) error {
	ctx := context.Background()
	sub, err := client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	syslog.Printf("Subscription creation err: %v", err)
	s, ok := status.FromError(err)
	if !ok || s.Code() != codes.AlreadyExists {
		if err != nil {
			syslog.Fatalf("Subscrbtion creation failed: %v", err)
		}
		return err
	}
	syslog.Printf("Created subscription: %v\n", sub)
	return nil
}
