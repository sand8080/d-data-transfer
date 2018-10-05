package main

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"path"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/sand8080/d-data-transfer/internal/data"
	"github.com/sand8080/d-data-transfer/internal/env"
	"github.com/sand8080/d-data-transfer/internal/handler"
	"github.com/sand8080/d-data-transfer/internal/validator"
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
		{UID: "unknown", IncomingCallNumber: "push"},
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

func AddEventsHandler(url string, sp *validator.SchemaProvider) *handler.JSONHandler {
	stdlog.Println("Adding events handler")
	cwd, _ := os.Getwd()
	schemaPath := path.Join(cwd, "schema/api/v1/events/add/request.json")
	stdlog.Printf("Loading schema from: %v\n", schemaPath)
	err := sp.Register(url, schemaPath)
	if err != nil {
		stdlog.Fatalf("Events handler adding failed: %v\n", err)
	}
	h, err := handler.NewJSONHandler(url, sp, addEventsProcessor)
	if err != nil {
		stdlog.Fatalf("AddEventsHanlder registration error: %v\n", err)
	}
	stdlog.Printf("Events handler added\n")
	return h
}
