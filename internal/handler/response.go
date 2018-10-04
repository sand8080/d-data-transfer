package handler

import (
	"encoding/json"
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

func (r Response) Bytes() []byte {
	if b, err := json.Marshal(r); err != nil {
		return []byte(err.Error())
	} else {
		return b
	}
}

func newError(code int, e error) *Response {
	return &Response{Code: code, Status: StatusError, Error: e.Error()}
}

func newErrorFromValidationErrors(code int, errs []gojsonschema.ResultError) *Response {
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.String())
	}
	message := strings.Join(msgs, "\r\n")

	return &Response{Code: code, Status: StatusError, Error: message}
}
