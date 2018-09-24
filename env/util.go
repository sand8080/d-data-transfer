package env

import "os"

func MustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		os.Exit(1)
	}
	return v
}

func Project() string {
	return MustGetenv("GOOGLE_CLOUD_PROJECT")
}

func Dataset() string {
	return MustGetenv("DATASET_ID")
}

func EventsTable() string {
	return MustGetenv("EVENTS_TABLE_ID")
}

func DataSourceUrl() string {
	return MustGetenv("DATA_SOURCE_URL")
}

func DataSourceAuthKey() string {
	return MustGetenv("DATA_SOURCE_AUTH_KEY")
}
