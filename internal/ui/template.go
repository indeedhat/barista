package ui

import (
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/indeedhat/barista/assets/templates"
)

func init() {
	tmpl = template.New("")

	tmpls = template.Must(
		tmpl.Funcs(templateFuncs).
			ParseFS(templates.FS, "layouts/*", "pages/*", "components/*"),
	)
}

var (
	tmpl  *template.Template
	tmpls *template.Template
)

var templateFuncs = template.FuncMap{
	"embed": func(name string, data any) template.HTML {
		var out strings.Builder
		if err := tmpl.ExecuteTemplate(&out, name, data); err != nil {
			log.Println(err)
		}
		return template.HTML(out.String())
	},
	"selected": func(actual, expect any) string {
		if fmt.Sprint(expect) == fmt.Sprint(actual) {
			return "selected"
		}
		return ""
	},
	"checked": func(actual, expect any) string {
		if fmt.Sprint(expect) == fmt.Sprint(actual) {
			return "checked"
		}
		return ""
	},
	"unix": func() int {
		return int(time.Now().Unix())
	},
	"map": func(values ...any) map[string]any {
		m := make(map[string]any)

		for i := 0; i < len(values); i += 2 {
			if i+1 >= len(values) {
				m[values[i].(string)] = nil
			} else {
				m[values[i].(string)] = values[i+1]
			}
		}

		return m
	},
	"bool": func(v bool) string {
		if v {
			return "true"
		}
		return "false"
	},
}
