package form

import (
	"net/http"
	"strconv"
	"strings"
)

type SettingsForm struct {
	Name                 string
	Css                  string
	Footer               string
	PerPage              int64
	URL                  string
	Lang                 string
	SelectedRenderEngine string
}

func NewSettingsForm(r *http.Request) *SettingsForm {
	perPage, _ := strconv.ParseInt(r.FormValue("per-page"), 10, 64)
	return &SettingsForm{
		Name:                 strings.TrimSpace(r.FormValue("name")),
		Css:                  r.FormValue("css"),
		Footer:               r.FormValue("footer"),
		URL:                  r.FormValue("url"),
		Lang:                 r.FormValue("lang"),
		PerPage:              perPage,
		SelectedRenderEngine: r.FormValue("rendering-engine"),
	}
}
