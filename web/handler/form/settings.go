package form

import "net/http"

type SettingsForm struct {
	Name string
	Css  string
}

func NewSettingsForm(r *http.Request) *SettingsForm {
	return &SettingsForm{
		Name: r.FormValue("name"),
		Css:  r.FormValue("css"),
	}
}
