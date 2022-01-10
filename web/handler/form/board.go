package form

import (
	"net/http"
)

type BoardForm struct {
	Name        string
	Description string
}

func NewBoardForm(r *http.Request) *BoardForm {
	return &BoardForm{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}
}
