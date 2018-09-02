package main

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine/log"
)

func main() {
	client := newPubSubClient()
	topic := createTopic(client)
	sub := mustGetenv("PUBSUB_SUBS")
	createSubs(client, sub, topic)

	ctx := context.Background()
	if err := pullMsgs(client, sub, topic); err != nil {
		log.Debugf(ctx, "Pull messages failed: %v", err)
	}
}

func pullMsgs(client *pubsub.Client, subName string, topic *pubsub.Topic) error {
	ctx := context.Background()
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subName)
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		log.Debugf(ctx, "Got message: %q\n", string(msg.Data))
		mu.Lock()
		defer mu.Unlock()
		received++
	})
	if err != nil {
		return err
	}
	return nil
}
