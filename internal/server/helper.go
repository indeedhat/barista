package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// Body reads in the Request body as a byte array
func Body(r *http.Request) []byte {
	cfg := r.Context().(Context).Config()

	data, err := io.ReadAll(io.LimitReader(r.Body, cfg.MaxBodySize))
	if err != nil {
		return nil
	}

	r.Body = io.NopCloser(bytes.NewBuffer(data))

	return data
}

// UnmarshalBody unmarshales the request body into the provided data structure
//
// NB: this is JSON only
func UnmarshalBody(r *http.Request, v any) error {
	data := Body(r)
	if data == nil {
		return errors.New("could not read request body")
	}

	return json.Unmarshal(data, v)
}

// ValidateRequest runs the provided struct against its validation tags
func ValidateRequest(v any) error {
	checker := validator.New()
	return checker.Struct(v)
}

type errorResponse struct {
	Error string `json:"error"`
}

// WriteResponse is a convenience method for constructing a response to return from a controller
func WriteResponse(rw http.ResponseWriter, code int, v any) {
	var resp []byte

	switch val := v.(type) {
	case string:
		resp = []byte(val)
	case validator.ValidationErrors:
		v = extractFIeldErrors(val)
	case error:
		v = errorResponse{val.Error()}
	}

	if resp == nil && v != nil {
		data, err := json.Marshal(v)
		if err != nil {
			code = http.StatusInternalServerError
			resp = []byte(`{"error":"failed to generate response json"}`)
		}

		resp = data
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)
	if v != nil {
		rw.Write(resp)
	}
}

type fieldErrorsResponse struct {
	Fields map[string][]string `json:"fields"`
}

func extractFIeldErrors(errs validator.ValidationErrors) fieldErrorsResponse {
	resp := fieldErrorsResponse{
		Fields: make(map[string][]string),
	}

	for _, err := range errs {
		k := err.Field()
		resp.Fields[k] = append(resp.Fields[k], err.Error())
	}

	return resp
}

func PathID(r *http.Request) (uint, error) {
	id := r.PathValue("id")
	if id == "" {
		return 0, errors.New("id not found in path")
	}

	n, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("id is not an int")
	}

	if n < 0 {
		return 0, errors.New("id must not be negative")
	}

	return uint(n), nil
}
