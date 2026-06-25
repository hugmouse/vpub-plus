package form

import (
	"errors"
	"net/http"
	"strings"
)

type SetupForm struct {
	ForumName     string
	URL           string
	Lang          string
	AdminName     string
	AdminPassword string
	Confirm       string
}

func (f *SetupForm) Validate() error {
	if strings.TrimSpace(f.ForumName) == "" {
		return errors.New("forum name is required")
	}
	if strings.TrimSpace(f.AdminName) == "" {
		return errors.New("admin username is required")
	}
	if f.AdminPassword == "" {
		return errors.New("admin password is required")
	}
	if f.AdminPassword != f.Confirm {
		return errors.New("password doesn't match confirmation")
	}
	return nil
}

func NewSetupForm(r *http.Request) *SetupForm {
	return &SetupForm{
		ForumName:     strings.TrimSpace(r.FormValue("forum-name")),
		URL:           strings.TrimSpace(r.FormValue("url")),
		Lang:          r.FormValue("lang"),
		AdminName:     strings.TrimSpace(r.FormValue("admin-name")),
		AdminPassword: r.FormValue("admin-password"),
		Confirm:       r.FormValue("confirm"),
	}
}
