package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

const StatusOK = "ok"
const StatusError = "error"

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func NewResponse(result interface{}) *Response {
	return &Response{
		Code:   http.StatusOK,
		Status: StatusOK,
		Result: result,
	}
}

func (r Response) Bytes() []byte {
	if b, err := json.Marshal(r); err != nil {
		return []byte(err.Error())
	} else {
		return b
	}
}

func NewErrorResponse(code int, e error) *Response {
	return &Response{Code: code, Status: StatusError, Error: e.Error()}
}

func NewErrorResponseFromValidationErrors(code int, errs []gojsonschema.ResultError) *Response {
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.String())
	}
	message := strings.Join(msgs, "\r\n")

	return &Response{Code: code, Status: StatusError, Error: message}
}
