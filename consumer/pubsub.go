package main

import (
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
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
