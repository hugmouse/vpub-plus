package form

import (
	"net/http"
	"strconv"
	"strings"
	"vpub/model"
)

type ThreadForm struct {
	Subject string
	Content string
	Topic   model.Topic
}

func NewThreadForm(r *http.Request) *ThreadForm {
	topicId, _ := strconv.ParseInt(r.FormValue("topicId"), 10, 64)
	return &ThreadForm{
		Subject: strings.TrimSpace(r.FormValue("subject")),
		Content: strings.TrimSpace(r.FormValue("content")),
		Topic: model.Topic{
			Id: topicId,
		},
	}
}
