package main

import (
	"net/http"

	"google.golang.org/appengine"

	"github.com/sand8080/d-data-transfer/internal/validator"
)

func main() {
	sp := validator.NewSchemaProvider()

	// Add events handler
	url := "/api/v1/events/add"
	h := AddEventsHandler(url, sp)
	http.HandleFunc(url, h.Handle)
	//
	//http.HandleFunc("/push", push)
	//http.HandleFunc("/initDataset", initDataset)
	appengine.Main()
}
