package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/indeedhat/barista/internal/ui"
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
func UnmarshalBody(r *http.Request, v any, pageData ...ui.Former) error {
	data := Body(r)
	if data == nil {
		return errors.New("could not read request body")
	}

	err := json.Unmarshal(data, v)

	if err == nil && len(pageData) > 0 {
		pageData[0].SetForm(v)
	}

	return err
}

// ValidateRequest runs the provided struct against its validation tags
func ValidateRequest(v any, pageData ...ui.ErrorFielder) error {
	checker := validator.New()
	checker.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	err := checker.Struct(v)
	if err != nil && len(pageData) > 0 {
		pageData[0].SetFieldErrors(ExtractFIeldErrors(err).Fields)
	}

	return err
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
		v = ExtractFIeldErrors(val)
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

func ExtractFIeldErrors(errs error) fieldErrorsResponse {
	resp := fieldErrorsResponse{
		Fields: make(map[string][]string),
	}

	for _, err := range errs.(validator.ValidationErrors) {
		k := err.Field()
		resp.Fields[k] = append(
			resp.Fields[k],
			fmt.Sprintf("validation for '%s' failed on the '%s' tag", k, err.Tag()),
		)
	}

	return resp
}

func PathID(r *http.Request, k ...string) (uint, error) {
	if len(k) == 0 {
		k = append(k, "id")
	}

	id := r.PathValue(k[0])
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

type UploadProps struct {
	Optional bool
	Ext      []string
	Mime     []string
}

func UploadFile(r *http.Request, formKey, savePath string, props *UploadProps) (string, error) {
	file, header, err := r.FormFile(formKey)
	if err != nil {
		if props != nil && props.Optional && errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		return "", errors.New("file upload failed")
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(header.Filename))
	if props != nil && len(props.Ext) > 0 {
		if !slices.Contains(props.Ext, ext) {
			return "", fmt.Errorf("file extension %s not allowed", ext)
		}
	}

	if props != nil && len(props.Mime) > 0 {
		buf := make([]byte, 512)
		if _, err := file.Read(buf); err != nil {
			return "", errors.New("filetype could not be verified")
		}

		mimeType := http.DetectContentType(buf)
		if !slices.Contains(props.Mime, mimeType) {
			return "", fmt.Errorf("mime type %s not allowed", mimeType)
		}
	}

	file.Seek(0, io.SeekStart)

	_ = os.MkdirAll(path.Dir(savePath), os.ModePerm)
	saveFile, err := os.Create(savePath + ext)
	if err != nil {
		return "", errors.New("file upload could not be saved on the server")
	}
	defer saveFile.Close()

	_, err = io.Copy(saveFile, file)
	if err != nil {
		return "", errors.New("file upload could not be saved on the server")
	}

	return savePath + ext, nil
}

func Redirect(rw http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(rw, r, url, http.StatusSeeOther)
}
