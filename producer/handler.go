package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/sand8080/d-data-transfer/env"
	"github.com/sand8080/d-data-transfer/queue"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func initQueue(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Initializing producer")

	// Creating pubsub client
	cli, err := queue.NewPubSubClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creating pubsub topic
	topic, err := queue.CreateTopic(ctx, cli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creating pubsub subscription
	_, err = queue.CreateSubs(ctx, cli, topic, env.MustGetenv("PUBSUB_SUBS"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Ok (init)")
	log.Debugf(ctx, "Producer initialization complete")
}

func publish(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Publishing new message")

	cli, err := queue.NewPubSubClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topic := queue.GetTopic(cli)

	payload := r.FormValue("payload")
	msg := pubsub.Message{
		Data: []byte(payload),
	}
	res := topic.Publish(ctx, &msg)
	_, err = res.Get(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Message is published")
	log.Debugf(ctx, "New message is published: %v", payload)
}
