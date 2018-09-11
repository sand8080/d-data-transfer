package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
	log.Debugf(ctx, "Message: %v", req)
	fmt.Fprintf(w, "ok")
}
