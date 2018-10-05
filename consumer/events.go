package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sand8080/d-data-transfer/internal/data"
	"github.com/sand8080/d-data-transfer/internal/handler"
)

type events struct {
	Events []data.Event `json:"events"`
}

func addEventsProcessor(b []byte) *handler.Response {
	var evs events
	err := json.Unmarshal(b, &evs)
	if err != nil {
		return handler.NewErrorResponse(http.StatusInternalServerError, err)
	}

	for _, e := range evs.Events {
		log.Printf("Event: %v\n", e)
	}

	return handler.NewResponse(evs)
}
