package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/init", initQueue)
	http.HandleFunc("/publish", publish)
	appengine.Main()
}
