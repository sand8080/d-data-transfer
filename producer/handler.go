package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type handler struct{}

func (h *handler) init(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Initializing producer")

	// Creating pubsub client
	cli, err := newPubSubClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creating pubsub topic
	topic, err := createTopic(ctx, cli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creating pubsub subscription
	_, err = createSubs(ctx, cli, topic, mustGetenv("PUBSUB_SUBS"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Ok (init)")
	log.Debugf(ctx, "Producer initialization complete")
}

func (h *handler) publish(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Publishing new message")

	cli, err := newPubSubClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topic := getTopic(cli)

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

func (h *handler) push(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Handling pushed message")
	fmt.Fprintf(w, "ok")
}
