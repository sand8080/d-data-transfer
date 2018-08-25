package main

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Printf("%s environment variable not set.", k)
		os.Exit(1)
	}
	return v
}

func newPubSubClient() *pubsub.Client {
	project := mustGetenv("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		fmt.Printf("Could not create pubsub Client: %v", err)
		os.Exit(1)
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
		fmt.Println(err)
		os.Exit(1)
	}
	if ok {
		return t
	}

	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		fmt.Printf("Failed to create the topic: %v", err)
		os.Exit(1)
	}
	return t
}

func createSubs(client *pubsub.Client, subName string, topic *pubsub.Topic) error {
	ctx := context.Background()
	sub, err := client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	fmt.Printf("Subscription creation err: %v", err)
	s, ok := status.FromError(err)
	if !ok || s.Code() != codes.AlreadyExists {
		if err != nil {
			fmt.Printf("Subscrbtion creation failed: %v", err)
			os.Exit(1)
		}
		return err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	return nil
}
