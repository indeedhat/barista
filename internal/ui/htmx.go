package ui

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Former interface {
	SetForm(v any)
}

type ErrorFielder interface {
	SetFieldErrors(v map[string][]string)
}

type ComponentData map[string]any

func NewComponentData(component string, d ...ComponentData) ComponentData {
	var data ComponentData
	if len(d) > 0 {
		data = d[0]
		data["Component"] = component
	} else {
		data = ComponentData{
			"Component": component,
		}
	}

	return data
}

// SetRequest implements Requester.
func (c ComponentData) SetForm(v any) {
	c["Form"] = v
}

// SetFieldErrors implements Fielder.
func (c ComponentData) SetFieldErrors(v map[string][]string) {
	c["FieldErrors"] = v
}

var _ Former = (*ComponentData)(nil)
var _ ErrorFielder = (*ComponentData)(nil)

func RenderComponent(w http.ResponseWriter, data ComponentData) error {
	var component string
	if c, found := data["Component"].(string); found && c != "" {
		component = c
	} else {
		return errors.New("Component not set")
	}

	return tmpls.ExecuteTemplate(w, component, data)
}

type PageData struct {
	Title       string
	Page        string
	User        any
	Form        any
	FieldErrors map[string][]string
	Data        map[string]any
}

func NewPageData(title, page string, user ...any) PageData {
	d := PageData{
		Title:       title,
		Page:        "pages/" + page,
		Form:        make(map[string]any),
		FieldErrors: make(map[string][]string),
		Data:        make(map[string]any),
	}

	if len(user) > 0 {
		d.User = user[0]
	}

	return d
}

// SetRequest implements Requester.
func (p *PageData) SetForm(v any) {
	p.Form = v
}

// SetFieldErrors implements Fielder.
func (p *PageData) SetFieldErrors(v map[string][]string) {
	p.FieldErrors = v
}

var _ Former = (*ComponentData)(nil)
var _ ErrorFielder = (*ComponentData)(nil)

func RenderGuest(w http.ResponseWriter, r *http.Request, data PageData) error {
	if r.Header.Get("HX-Request") == "true" {
		return tmpls.ExecuteTemplate(w, "layouts/hx", data)
	}

	return tmpls.ExecuteTemplate(w, "layouts/guest", data)
}

func RenderUser(w http.ResponseWriter, r *http.Request, data PageData) error {
	if r.Header.Get("HX-Request") == "true" {
		return tmpls.ExecuteTemplate(w, "layouts/hx", data)
	}

	return tmpls.ExecuteTemplate(w, "layouts/user", data)
}

type ToastLevel string

const (
	Info    ToastLevel = "info"
	Warning ToastLevel = "error"
	Success ToastLevel = "success"
)

type toastEvent struct {
	Level   ToastLevel `json:"level"`
	Message string     `json:"message"`
}

func Toast(w http.ResponseWriter, level ToastLevel, message string, code ...int) error {
	data, err := json.Marshal(map[string]any{
		"triggerToast": toastEvent{level, message},
	})
	if err != nil {
		return err
	}

	w.Header().Set("HX-Trigger", string(data))
	if len(code) > 1 {
		w.WriteHeader(code[0])
	}

	return nil
}

func Redirect(rw http.ResponseWriter, url string) {
	rw.Header().Set("HX-Redirect", url)
	rw.WriteHeader(http.StatusNoContent)
}
