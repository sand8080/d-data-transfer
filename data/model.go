package data

import "cloud.google.com/go/bigquery"

type Event struct {
	ID     string
	Source string
	Created Datetime
}

// Save implements the ValueSaver interface.
func (e *Event) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"ID":     e.ID,
		"Source": e.Source,
	}, "", nil
}
