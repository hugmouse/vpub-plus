package handler

import (
	"fmt"
	"html/template"
	"io"
	"time"
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
			"iso8601": func(t time.Time) string {
				return t.Format("2006-01-02")
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
		//data["hasNotifications"] = h.storage.UserHasNotifications(user)
	}
	data["logged"] = user
	data["boardTitle"] = h.title
	if err := views[view].Funcs(template.FuncMap{
		"hasPermission": func(name string) bool {
			return user == name
		},
		"logged": func() bool {
			return user != ""
		},
	}).ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Println(err)
	}

}
