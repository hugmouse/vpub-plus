package form

import "net/http"

type GroupForm struct {
	Name string
}

func NewGroupForm(r *http.Request) *GroupForm {
	return &GroupForm{
		Name: r.FormValue("name"),
	}
}
