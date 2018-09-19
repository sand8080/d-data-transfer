package main

import (
	"google.golang.org/appengine"
	"net/http"
)

func main() {
	http.HandleFunc("/push", push)
	http.HandleFunc("/initDataset", initDataset)
	appengine.Main()
}
