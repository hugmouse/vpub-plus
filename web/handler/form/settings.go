package form

import (
	"net/http"
	"strconv"
)

type SettingsForm struct {
	Name    string
	Css     string
	Footer  string
	PerPage int64
}

func NewSettingsForm(r *http.Request) *SettingsForm {
	perPage, _ := strconv.ParseInt(r.FormValue("per-page"), 10, 64)
	return &SettingsForm{
		Name:    r.FormValue("name"),
		Css:     r.FormValue("css"),
		Footer:  r.FormValue("footer"),
		PerPage: perPage,
	}
}
