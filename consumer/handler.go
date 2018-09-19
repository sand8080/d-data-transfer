package main

import (
	"encoding/json"
	"fmt"
	"github.com/sand8080/d-data-transfer/env"
	"net/http"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/sand8080/d-data-transfer/data"
)

type pushRequest struct {
	Message      pubsub.Message
	Subscription string
}

func push(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Handling pushed message")
	var req pushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cli, err := data.NewBQClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := cli.Dataset(env.Dataset()).Table(env.EventsTable()).Uploader()
	events := []data.Event{
		{ID: "unknown", Source: "push"},
	}
	err = u.Put(ctx, events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf(ctx, "Message: %v", req)
	fmt.Fprintf(w, "ok")
}

func initDataset(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "Dataset initialization started")

	cli, err := data.NewBQClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = data.CreateEventsTable(ctx, cli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf(ctx, "Dataset initialized")
	fmt.Fprintf(w, "Dataset initialized")
}

