package validator

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type SchemaProvider struct {
	validators map[string]*gojsonschema.Schema
}

func NewSchemaProvider() *SchemaProvider {
	p := SchemaProvider{validators: make(map[string]*gojsonschema.Schema)}
	return &p
}

func loadSchema(path string) (*gojsonschema.Schema, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Schema abs path error: %v\n", path)
		return nil, err
	}

	path = filepath.ToSlash(path)
	path = fmt.Sprintf("file://%s", path)

	loader := gojsonschema.NewReferenceLoader(path)
	return gojsonschema.NewSchema(loader)
}

func prepareUrl(url string) string {
	return strings.ToLower(strings.Trim(url, "/"))
}

func (p *SchemaProvider) Register(url, filename string) error {
	schema, err := loadSchema(filename)
	if err != nil {
		log.Printf("Schema loading error: %v\n", err)
		return err
	}

	p.validators[prepareUrl(url)] = schema
	return nil
}

func (p *SchemaProvider) Get(url string) (*gojsonschema.Schema, error) {
	s, ok := p.validators[prepareUrl(url)]
	if ok {
		return s, nil
	} else {
		return nil, errors.New("schema not found")
	}
}
