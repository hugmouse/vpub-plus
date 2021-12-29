package form

import (
	"net/http"
)

type UserForm struct {
	Username string
	Password string
	Key      string
}

func NewUserForm(r *http.Request) *UserForm {
	return &UserForm{
		Username: r.FormValue("name"),
		Password: r.FormValue("password"),
		Key:      r.FormValue("key"),
	}
}
