package handler

import (
	"html/template"
	"io"
	"vpub/gmi2html"
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
			"logged": func() bool {
				return false
			},
			"gmi2html": func(gmi string) template.HTML {
				return template.HTML(gmi2html.Convert(gmi))
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
	if user != "" {
		data["hasNotifications"] = h.storage.UserHasNotifications(user)
	}
	data["logged"] = user
	data["boardTitle"] = h.title
	views[view].Funcs(template.FuncMap{
		"hasPermission": func(name string) bool {
			return user == name
		},
		"logged": func() bool {
			return user != ""
		},
	}).ExecuteTemplate(w, "layout", data)
}

func (h *Handler) view(view string) *template.Template {
	return views[view]
}
