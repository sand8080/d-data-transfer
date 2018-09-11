package main

import (
	"context"
	"net/http"
	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func main() {
	http.HandleFunc("/push", push)
	appengine.Main()
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
