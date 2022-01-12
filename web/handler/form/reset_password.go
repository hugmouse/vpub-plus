package form

import (
	"errors"
	"net/http"
)

type ResetPasswordForm struct {
	Password string
	Confirm  string
	Hash     string
}

func (f *ResetPasswordForm) Validate() error {
	if f.Password != f.Confirm {
		return errors.New("password doesn't match confirmation")
	}
	if len(f.Password) < 6 {
		return errors.New("password needs to be at least 6 characters")
	}
	return nil
}

func NewResetPasswordForm(r *http.Request) *ResetPasswordForm {
	return &ResetPasswordForm{
		Confirm:  r.FormValue("confirm"),
		Password: r.FormValue("password"),
		Hash:     r.FormValue("hash"),
	}
}
