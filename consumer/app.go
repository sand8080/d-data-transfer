package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

func main() {
	client := newPubSubClient()
	topic := createTopic(client)

	if err := pullMsgs(client, mustGetenv("PUBSUB_SUBS"), topic); err != nil {
		log.Fatal(err)
	}
}

func pullMsgs(client *pubsub.Client, subName string, topic *pubsub.Topic) error {
	ctx := context.Background()
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subName)
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Got message: %q\n", string(msg.Data))
		mu.Lock()
		defer mu.Unlock()
		received++
	})
	if err != nil {
		return err
	}
	return nil
}
