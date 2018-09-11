package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", root)
	appengine.Main()
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default service app stub")
}
