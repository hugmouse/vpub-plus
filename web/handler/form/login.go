package form

import (
	"net/http"
)

type LoginForm struct {
	Username string
	Password string
}

func NewLoginForm(r *http.Request) *LoginForm {
	return &LoginForm{
		Username: r.FormValue("name"),
		Password: r.FormValue("password"),
	}
}
