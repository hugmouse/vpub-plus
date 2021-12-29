package form

import "net/http"

type ReplyForm struct {
	Content string
}

func NewReplyForm(r *http.Request) *ReplyForm {
	return &ReplyForm{
		Content: r.FormValue("reply"),
	}
}
