package form

import "net/http"

type AccountForm struct {
	Picture string
	About   string
}

func NewAccountForm(r *http.Request) *AccountForm {
	return &AccountForm{
		Picture: r.FormValue("picture"),
		About:   r.FormValue("about"),
	}
}
