package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	h := handler{}
	http.HandleFunc("/init", h.init)
	http.HandleFunc("/publish", h.publish)
	http.HandleFunc("/push", h.push)
	appengine.Main()
}
