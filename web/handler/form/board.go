package form

import (
	"net/http"
	"strconv"
	"vpub/model"
)

type BoardForm struct {
	Name        string
	Description string
	Position    int64
	Forums      []model.Forum
	ForumId     int64
	IsLocked    bool
}

func NewBoardForm(r *http.Request) *BoardForm {
	position, _ := strconv.ParseInt(r.FormValue("position"), 10, 64)
	forumId, _ := strconv.ParseInt(r.FormValue("forumId"), 10, 64)
	return &BoardForm{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Position:    position,
		ForumId:     forumId,
		IsLocked:    r.FormValue("locked") == "true",
	}
}
