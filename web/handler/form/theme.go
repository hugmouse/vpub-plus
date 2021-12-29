package form

import (
	"net/http"
)

type ThemeForm struct {
	Theme string
}

func NewThemeForm(r *http.Request) *ThemeForm {
	return &ThemeForm{
		Theme: r.FormValue("theme"),
	}
}
