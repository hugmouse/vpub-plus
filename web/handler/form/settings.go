package form

import "net/http"

type SettingsForm struct {
	Name   string
	Css    string
	Footer string
}

func NewSettingsForm(r *http.Request) *SettingsForm {
	return &SettingsForm{
		Name:   r.FormValue("name"),
		Css:    r.FormValue("css"),
		Footer: r.FormValue("footer"),
	}
}
