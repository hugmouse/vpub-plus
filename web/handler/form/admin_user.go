package form

import (
	"net/http"
)

type AdminUserForm struct {
	Username string
	About    string
}

func NewAdminUserForm(r *http.Request) *AdminUserForm {
	return &AdminUserForm{
		Username: r.FormValue("name"),
		About:    r.FormValue("about"),
	}
}
