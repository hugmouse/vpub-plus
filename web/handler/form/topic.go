package form

import (
	"net/http"
	"strconv"
	"vpub/model"
)

type TopicForm struct {
	ID         int64
	BoardID    int64
	PostForm   PostForm
	IsSticky   bool
	IsLocked   bool
	Boards     []model.Board
	NewBoardID int64
}

func NewTopicForm(r *http.Request) *TopicForm {
	BoardID, _ := strconv.ParseInt(r.FormValue("boardId"), 10, 64)
	NewBoardID, _ := strconv.ParseInt(r.FormValue("newBoardId"), 10, 64)
	return &TopicForm{
		BoardID:    BoardID,
		PostForm:   NewPostForm(r),
		IsSticky:   r.FormValue("sticky") == "on",
		IsLocked:   r.FormValue("locked") == "on",
		NewBoardID: NewBoardID,
	}
}
