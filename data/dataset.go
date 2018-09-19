package data

import (
	"context"

	"cloud.google.com/go/bigquery"

	"github.com/sand8080/d-data-transfer/env"
)

func NewBQClient(ctx context.Context) (*bigquery.Client, error) {
	return bigquery.NewClient(ctx, env.Project())
}

func CreateEventsTable(ctx context.Context, client *bigquery.Client) error {
	// bigquery.InferSchema infers BQ schema from native Go types.
	schema, err := bigquery.InferSchema(Event{})
	if err != nil {
		return err
	}

	table := client.Dataset(env.Dataset()).Table(env.EventsTable())
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		return err
	}
	return nil
}