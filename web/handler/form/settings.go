package form

import (
	"net/http"
	"strconv"
	"strings"
	"vpub/model"
)

type SettingsForm struct {
	Name    string
	Css     string
	Footer  string
	PerPage int64
}

func (f *SettingsForm) Merge(settings *model.Settings) *model.Settings {
	settings.Name = f.Name
	settings.Css = f.Css
	settings.Footer = f.Footer
	settings.PerPage = f.PerPage
	return settings
}

func NewSettingsForm(r *http.Request) *SettingsForm {
	perPage, _ := strconv.ParseInt(r.FormValue("per-page"), 10, 64)
	return &SettingsForm{
		Name:    strings.TrimSpace(r.FormValue("name")),
		Css:     r.FormValue("css"),
		Footer:  r.FormValue("footer"),
		PerPage: perPage,
	}
}
