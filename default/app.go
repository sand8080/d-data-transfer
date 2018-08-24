package main

import (
	"net/http"

	"fmt"
	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", addCronTask)
	appengine.Main()
}

func addCronTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default service app stub")
}
