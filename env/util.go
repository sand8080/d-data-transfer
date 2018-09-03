package env

import "os"

func MustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		os.Exit(1)
	}
	return v
}
