package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sand8080/d-data-transfer/internal/validator"
)

func createTmpSchema(schema []byte) (*os.File, error) {
	f, err := ioutil.TempFile(os.TempDir(), "d-data-transfer-schema-")
	if err != nil {
		return nil, err
	}
	_, err = f.Write(schema)
	if err != nil {
		os.Remove(f.Name())
		return nil, err
	}
	f.Close()

	return f, nil
}

func prepareSchemaProvider(t *testing.T, url string) *validator.SchemaProvider {
	schema := []byte(`{
		"$schema": "http://json-schema.org/schema#",
		"type": "object",
		"properties": {
			"uid": {"type": "string"}
		},
		"required": ["uid"],
		"additionalProperties": false
	}`)
	tmpFile, err := createTmpSchema(schema)
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	p := validator.NewSchemaProvider()
	p.Register(url, tmpFile.Name())

	return p
}

func TestJSONHandler_Handle(t *testing.T) {
	url := "/do"
	p := prepareSchemaProvider(t, url)
	proc := func(data []byte) *Response {
		return &Response{Code: http.StatusOK, Status: StatusOK, Result: string(data)}
	}
	h, err := NewJSONHandler(url, p, proc)
	assert.NoError(t, err)

	cases := []struct {
		reqData   []byte
		expCode   int
		expStatus string
		expError  string
	}{
		{[]byte(`{"uid": "DtM56"}`), http.StatusOK, StatusOK, ""},
		{[]byte(`{"uid": 42}`), http.StatusBadRequest, StatusError,
			"uid: Invalid type. Expected: string, given: integer"},
		{[]byte(``), http.StatusBadRequest, StatusError, "empty request"},
	}
	for _, c := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/do", bytes.NewReader(c.reqData))
		h.Handle(w, req)

		resp := w.Result()
		assert.Equal(t, c.expCode, resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)

		var r Response
		json.Unmarshal(body, &r)
		assert.Equal(t, c.expCode, r.Code)
		assert.Equal(t, c.expStatus, r.Status)
		if c.expCode == http.StatusOK {
			assert.Nil(t, r.Error)
			assert.NotNil(t, r.Result)
		} else {
			assert.NotNil(t, r.Error)
			assert.Nil(t, r.Result)
			assert.Equal(t, c.expError, r.Error)
		}
		fmt.Printf("Response: %v\n", r)
	}
}
