package form

import (
	"errors"
	"net/http"
	"regexp"
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
	if len(f.Username) < 3 {
		return errors.New("username needs to be at least 3 characters")
	}
	if len(f.Username) > 20 {
		return errors.New("username should be 20 characters or less")
	}
	match, _ := regexp.MatchString("^[a-z0-9-_]+$", f.Username)
	if !match {
		return errors.New("only lowercase letters and digits are accepted for username")
	}
	if len(f.Password) < 6 {
		return errors.New("password needs to be at least 6 characters")
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
