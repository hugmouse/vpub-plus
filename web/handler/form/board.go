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
	ForumID     int64
	IsLocked    bool
}

func (f *BoardForm) Merge(board *model.Board) *model.Board {
	board.Name = f.Name
	board.Description = f.Description
	board.Position = f.Position
	board.IsLocked = f.IsLocked
	board.Forum = model.Forum{ID: f.ForumID}
	return board
}

func NewBoardForm(r *http.Request) *BoardForm {
	position, _ := strconv.ParseInt(r.FormValue("position"), 10, 64)
	forumId, _ := strconv.ParseInt(r.FormValue("forumId"), 10, 64)
	return &BoardForm{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Position:    position,
		ForumID:     forumId,
		IsLocked:    r.FormValue("locked") == "true",
	}
}
