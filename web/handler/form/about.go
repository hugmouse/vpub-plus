package form

import "net/http"

type AboutForm struct {
	About string
}

func NewAboutForm(r *http.Request) *AboutForm {
	return &AboutForm{
		About: r.FormValue("about"),
	}
}
