package handler

import (
	"html/template"
	"io"
)

var views = make(map[string]*template.Template)

func (h *Handler) initTpl() {
	commonTemplates := ""
	for _, content := range TplCommonMap {
		commonTemplates += content
	}

	for name, content := range TplMap {
		views[name] = template.Must(template.New("main").Funcs(template.FuncMap{
			"hasPermission": func(name string) bool {
				return false
			},
		}).Parse(commonTemplates + content))
	}
}

func (h *Handler) renderLayout(w io.Writer, view string, params map[string]interface{}, user string) {
	data := make(map[string]interface{})
	if params != nil {
		for k, v := range params {
			data[k] = v
		}
	}
	data["logged"] = user
	data["boardTitle"] = h.title
	views[view].Funcs(template.FuncMap{
		"hasPermission": func(name string) bool {
			return user == name
		},
	}).ExecuteTemplate(w, "layout", data)
}
