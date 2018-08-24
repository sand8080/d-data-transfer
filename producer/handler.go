package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type handler struct {
	topic *pubsub.Topic
}

func NewHandler(topic *pubsub.Topic) *handler {
	return &handler{topic: topic}
}

func (h *handler) addCronTask(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Scheduling new cron task")

	msg := pubsub.Message{
		Data: []byte("cron task"),
	}
	res := h.topic.Publish(ctx, &msg)
	_, err := res.Get(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Message is published")
	log.Debugf(ctx, "New cron task scheduled")
}
