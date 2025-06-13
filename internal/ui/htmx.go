package ui

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/indeedhat/barista/assets/templates"
)

func init() {
	tmpl := template.New("")
	tmpls = template.Must(tmpl.Funcs(template.FuncMap{
		"embed": func(name string, data any) template.HTML {
			var out strings.Builder
			if err := tmpl.ExecuteTemplate(&out, name, data); err != nil {
				log.Println(err)
			}
			return template.HTML(out.String())
		},
	}).ParseFS(templates.FS, "layouts/*", "pages/*", "components/*"))
}

var tmpls *template.Template

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
