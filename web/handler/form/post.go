package form

import (
	"net/http"
	"strconv"
	"strings"
)

type PostForm struct {
	Subject string
	Content string
	TopicID int64
}

func NewPostForm(r *http.Request) PostForm {
	TopicID, _ := strconv.ParseInt(r.FormValue("topicId"), 10, 64)
	return PostForm{
		Subject: strings.TrimSpace(r.FormValue("subject")),
		Content: r.FormValue("content"),
		TopicID: TopicID,
	}
}
