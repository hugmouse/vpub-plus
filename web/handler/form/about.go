package form

import "net/http"

type AccountForm struct {
	Picture    string
	PictureAlt string
	About      string
}

func NewAccountForm(r *http.Request) *AccountForm {
	return &AccountForm{
		Picture:    r.FormValue("picture"),
		PictureAlt: r.FormValue("picture-alt"),
		About:      r.FormValue("about"),
	}
}
