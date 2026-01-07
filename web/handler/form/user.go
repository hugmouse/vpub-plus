package form

import (
	"errors"
	"net/http"
	"vpub/model"
)

type UserForm struct {
	Username string
	Password string
	Confirm  string
	Key      string
}

func (f *UserForm) Validate() error {
	if f.Password != f.Confirm {
		return errors.New("password doesn't match confirmation")
	}

	return nil
}

func (f *UserForm) Merge(u *model.User) *model.User {
	u.Name = f.Username
	u.Password = f.Password
	return u
}

func NewUserForm(r *http.Request) *UserForm {
	return &UserForm{
		Username: r.FormValue("name"),
		Confirm:  r.FormValue("confirm"),
		Password: r.FormValue("password"),
		Key:      r.FormValue("key"),
	}
}
