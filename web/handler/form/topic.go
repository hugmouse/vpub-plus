package form

import (
	"net/http"
	"strconv"
)

type TopicForm struct {
	BoardId  int64
	PostForm PostForm
}

func NewTopicForm(r *http.Request) *TopicForm {
	BoardId, _ := strconv.ParseInt(r.FormValue("boardId"), 10, 64)
	return &TopicForm{
		BoardId:  BoardId,
		PostForm: NewPostForm(r),
	}
}
