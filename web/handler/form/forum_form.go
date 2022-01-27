package form

import (
	"net/http"
	"strconv"
)

type ForumForm struct {
	Name     string
	IsLocked bool
	Position int64
}

func NewForumForm(r *http.Request) *ForumForm {
	position, _ := strconv.ParseInt(r.FormValue("position"), 10, 64)
	return &ForumForm{
		Name:     r.FormValue("name"),
		IsLocked: r.FormValue("locked") == "on",
		Position: position,
	}
}
