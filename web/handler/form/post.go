package form

import (
	"net/http"
)

type PostForm struct {
	Title   string
	Content string
	Topics  []string
	Topic   string
}

func NewPostForm(r *http.Request, topics []string) *PostForm {
	return &PostForm{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Topic:   r.FormValue("topic"),
		Topics:  topics,
	}
}
