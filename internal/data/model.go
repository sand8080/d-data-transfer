package data

import "cloud.google.com/go/bigquery"

type Event struct {
	UID                string `json:"uid"`
	IncomingCallNumber string `json:"incomingCallNumber"`
}

// Save implements the ValueSaver interface.
func (e *Event) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"UID":                e.UID,
		"IncomingCallNumber": e.IncomingCallNumber,
	}, "", nil
}
