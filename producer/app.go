package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	client := newPubSubClient()
	topic := createTopic(client)
	createSubs(client, mustGetenv("PUBSUB_SUBS"), topic)
	h := handler{topic: topic}

	http.HandleFunc("/add-cron-task", h.addCronTask)
	appengine.Main()
}
