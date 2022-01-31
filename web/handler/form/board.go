package form

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
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

func (f *BoardForm) Validate() error {
	if len(strings.TrimSpace(f.Name)) == 0 {
		return errors.New("Board name can't be empty")
	}
	if len(strings.TrimSpace(f.Description)) == 0 {
		return errors.New("Description name can't be empty")
	}
	return nil
}

func (f *BoardForm) Merge(board *model.Board) *model.Board {
	board.Name = f.Name
	board.Description = f.Description
	board.Position = f.Position
	board.IsLocked = f.IsLocked
	board.Forum = model.Forum{Id: f.ForumId}
	return board
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
