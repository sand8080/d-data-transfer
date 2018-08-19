package main

import (
	"net/http"

	"google.golang.org/appengine"
	"fmt"
)

func main() {
	http.HandleFunc("/", addCronTask)
	appengine.Main()
}

func addCronTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default service app stub")
}
