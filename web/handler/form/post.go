package form

import (
	"net/http"
)

type PostForm struct {
	Title   string
	Content string
}

func NewPostForm(r *http.Request) *PostForm {
	return &PostForm{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}
}
