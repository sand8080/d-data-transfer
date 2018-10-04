package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"google.golang.org/appengine"

	"github.com/sand8080/d-data-transfer/internal/validator"
)

type H struct{}

func (h *H) handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok handle")
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Getting workdir error: %v\n", err)
	}

	p := validator.NewSchemaProvider()
	err = p.Register("/api/v1/event/add", path.Join(cwd, "schema/api/v1/event/add/request.json"))
	if err != nil {
		log.Fatalf("Registration error: %v\n", err)
	}

	h := H{}
	http.HandleFunc("/push", push)
	http.HandleFunc("/initDataset", initDataset)
	http.HandleFunc("/h", h.handle)
	appengine.Main()
}
