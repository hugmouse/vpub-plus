package form

import (
	"net/http"
	"strconv"
)

type BoardForm struct {
	Name        string
	Description string
	Position    int64
}

func NewBoardForm(r *http.Request) *BoardForm {
	position, _ := strconv.ParseInt(r.FormValue("position"), 10, 64)
	return &BoardForm{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Position:    position,
	}
}
