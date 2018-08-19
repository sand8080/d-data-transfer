package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"fmt"
)

func main() {
	http.HandleFunc("/add-cron-task", addCronTask)
	appengine.Main()
}

func addCronTask(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	log.Debugf(c, "Scheduling new cron task")
	fmt.Fprintln(w, "Hello, world!")
}
