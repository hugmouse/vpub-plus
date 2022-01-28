package handler

import (
	"fmt"
	"html/template"
	"io"
	"time"
	"vpub/model"
	"vpub/syntax"
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
			"syntax": func(input string) template.HTML {
				return template.HTML(syntax.Convert(input))
			},
			"iso8601": func(t time.Time) string {
				return t.Format("2006-01-02")
			},
			"iso8601Time": func(t time.Time) string {
				return t.Format("2006-01-02 15:04:05")
			},
			"html": func(s string) template.HTML {
				return template.HTML(s)
			},
			"timeAgo": func(t time.Time) string {
				d := time.Since(t)
				if d.Seconds() < 60 {
					seconds := int(d.Seconds())
					if seconds == 1 {
						return "1 second ago"
					}
					return fmt.Sprintf("%d seconds ago", seconds)
				} else if d.Minutes() < 60 {
					minutes := int(d.Minutes())
					if minutes == 1 {
						return "1 minute ago"
					}
					return fmt.Sprintf("%d minutes ago", minutes)
				} else if d.Hours() < 24 {
					hours := int(d.Hours())
					if hours == 1 {
						return "1 hour ago"
					}
					return fmt.Sprintf("%d hours ago", hours)
				} else {
					days := int(d.Hours()) / 24
					if days == 1 {
						return "1 day ago"
					}
					return fmt.Sprintf("%d days ago", days)
				}
			},
			"inc": func(v int64) int64 {
				return v + 1
			},
			"dec": func(v int64) int64 {
				return v - 1
			},
		}).Parse(commonTemplates + content))
	}
}

func (h *Handler) renderLayout(w io.Writer, view string, params map[string]interface{}, user model.User) {
	data := make(map[string]interface{})
	if params != nil {
		for k, v := range params {
			data[k] = v
		}
	}
	data["logged"] = user
	settings, err := h.storage.Settings()
	if err != nil {
		fmt.Println(err)
	}
	data["settings"] = settings
	if err := views[view].Funcs(template.FuncMap{
		"hasPermission": func(name string) bool {
			return user.Name == name
		},
		"logged": func() bool {
			return user.Name != ""
		},
	}).ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Println(err)
	}

}
