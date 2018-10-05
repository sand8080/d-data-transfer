package handler

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/xeipuuv/gojsonschema"

	"github.com/sand8080/d-data-transfer/internal/validator"
)

type Processor func(data []byte) *Response

type JSONHandler struct {
	url       string
	schema    *gojsonschema.Schema
	processor Processor
}

func NewJSONHandler(url string, p *validator.SchemaProvider, proc Processor) (*JSONHandler, error) {
	s, err := p.Get(url)
	if err != nil {
		return nil, err
	}
	return &JSONHandler{url: url, schema: s, processor: proc}, nil
}

func (h JSONHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Reading raw reqData
	dataLoader, reader := gojsonschema.NewReaderLoader(r.Body)
	reqData, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("Request data loading error: %v\n", err)
		write(w, NewErrorResponse(http.StatusInternalServerError, err))
		return
	}

	// Request validation
	valResult, err := h.schema.Validate(dataLoader)
	if err != nil {
		log.Printf("Request validation error: %v\n", err)
		if err == io.EOF {
			err = errors.New("empty request")
		}
		write(w, NewErrorResponse(http.StatusBadRequest, err))
		return
	}
	if !valResult.Valid() {
		log.Printf("Request invalid\n")
		write(w, NewErrorResponseFromValidationErrors(http.StatusBadRequest, valResult.Errors()))
		return
	}

	// Request processing
	result := h.processor(reqData)

	// TODO implement response validation

	// Writing response
	write(w, result)
	return
}

func write(w http.ResponseWriter, r *Response) {
	w.WriteHeader(r.Code)
	w.Write(r.Bytes())
}
