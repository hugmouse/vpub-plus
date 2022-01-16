package form

import (
	"net/http"
	"strconv"
	"strings"
)

type PostForm struct {
	Subject  string
	Content  string
	TopicId  int64
	BoardId  int64
	IsAdmin  bool
	IsSticky bool
}

func NewPostForm(r *http.Request) PostForm {
	TopicId, _ := strconv.ParseInt(r.FormValue("topicId"), 10, 64)
	BoardId, _ := strconv.ParseInt(r.FormValue("boardId"), 10, 64)
	return PostForm{
		Subject:  strings.TrimSpace(r.FormValue("subject")),
		Content:  r.FormValue("content"),
		TopicId:  TopicId,
		BoardId:  BoardId,
		IsSticky: r.FormValue("sticky") == "true",
	}
}
