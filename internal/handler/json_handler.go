package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/xeipuuv/gojsonschema"

	"github.com/sand8080/d-data-transfer/internal/validator"
)

type Processor func(data []byte) Response

type JSONHandler struct {
	schema    *gojsonschema.Schema
	processor Processor
}

func NewJSONHandler(url string, p *validator.SchemaProvider, proc Processor) (*JSONHandler, error) {
	s, err := p.Get(url)
	if err != nil {
		return nil, err
	}
	return &JSONHandler{s, proc}, nil
}

func (h JSONHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Reading raw reqData
	dataLoader, reader := gojsonschema.NewReaderLoader(r.Body)
	reqData, err := ioutil.ReadAll(reader)
	if err != nil {
		write(w, newError(http.StatusInternalServerError, err))
		return
	}

	// Request validation
	valResult, err := h.schema.Validate(dataLoader)
	if err != nil {
		write(w, newError(http.StatusInternalServerError, err))
		return
	}
	if !valResult.Valid() {
		write(w, newErrorFromValidationErrors(http.StatusBadRequest, valResult.Errors()))
		return
	}

	// Request processing
	result := h.processor(reqData)

	// TODO implement response validation

	// Writing response
	write(w, &result)
	return
}

func write(w http.ResponseWriter, r *Response) {
	w.WriteHeader(r.Code)
	w.Write(r.Bytes())
}
