package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	"time"
	"vpub/syntax"
	"vpub/web/handler/request"
)

var views = make(map[string]*template.Template)

type View struct {
	w      http.ResponseWriter
	r      *http.Request
	tpl    string
	params map[string]interface{}
}

func NewView(w http.ResponseWriter, r *http.Request, tpl string) View {
	params := make(map[string]interface{})
	params[csrf.TemplateTag] = csrf.TemplateField(r)
	return View{
		w:      w,
		r:      r,
		tpl:    tpl,
		params: params,
	}
}

func (v View) Set(key string, val interface{}) {
	v.params[key] = val
}

func (v View) Render() {
	user := request.GetUserContextKey(v.r)
	data := v.params
	data["logged"] = user
	settings := request.GetSettingsContextKey(v.r)
	data["settings"] = settings
	session := request.GetSessionContextKey(v.r)
	data["errors"] = session.GetFlashErrors()
	data["info"] = session.GetFlashInfo()
	session.Save(v.r, v.w)
	if err := views[v.tpl].Funcs(template.FuncMap{
		"hasPermission": func(name string) bool {
			return user.Name == name
		},
		"logged": func() bool {
			return user.Name != ""
		},
	}).ExecuteTemplate(v.w, "layout", data); err != nil {
		fmt.Println(err)
	}
}

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
				return template.HTML(syntax.Convert(input, true))
			},
			"sig": func(input string) template.HTML {
				return template.HTML(syntax.Convert(input, false))
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
				} else if d.Hours() < 730 {
					days := int(d.Hours()) / 24
					if days == 1 {
						return "1 day ago"
					}
					return fmt.Sprintf("%d days ago", days)
				} else if d.Hours() < 8760 {
					months := int(d.Hours()) / 730
					if months == 1 {
						return "1 month ago"
					}
					return fmt.Sprintf("%d months ago", months)
				} else {
					years := int(d.Hours()) / 8760
					if years == 1 {
						return "1 year ago"
					}
					return fmt.Sprintf("%d years ago", years)
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

func (h *Handler) renderLayout(w http.ResponseWriter, r *http.Request, view string, params map[string]interface{}) {
	user := request.GetUserContextKey(r)
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
	session, err := h.session.GetSession(r)
	if err != nil {
		fmt.Println(err)
	}
	data["errors"] = session.GetFlashErrors()
	data["info"] = session.GetFlashInfo()
	session.Save(r, w)
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
