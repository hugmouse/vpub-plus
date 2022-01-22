package form

import (
	"net/http"
	"strconv"
	"vpub/model"
)

type TopicForm struct {
	Id         int64
	BoardId    int64
	PostForm   PostForm
	IsSticky   bool
	IsLocked   bool
	Boards     []model.Board
	NewBoardId int64
}

func NewTopicForm(r *http.Request) *TopicForm {
	BoardId, _ := strconv.ParseInt(r.FormValue("boardId"), 10, 64)
	NewBoardId, _ := strconv.ParseInt(r.FormValue("newBoardId"), 10, 64)
	return &TopicForm{
		BoardId:    BoardId,
		PostForm:   NewPostForm(r),
		IsSticky:   r.FormValue("sticky") == "true",
		IsLocked:   r.FormValue("locked") == "true",
		NewBoardId: NewBoardId,
	}
}
