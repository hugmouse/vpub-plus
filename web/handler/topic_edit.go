package handler

import (
	"net/http"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) showEditTopicView(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)
	if !user.IsAdmin {
		notFound(w)
		return
	}
	topic, err := h.storage.TopicByID(RouteInt64Param(r, "topicId"))
	if err != nil {
		notFound(w)
		return
	}
	post, err := h.storage.PostByID(topic.Post.ID)
	if err != nil {
		notFound(w)
		return
	}
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	topicForm := form.TopicForm{
		ID:       topic.ID,
		BoardID:  topic.BoardID,
		IsSticky: topic.IsSticky,
		IsLocked: topic.IsLocked,
		Boards:   boards,
		PostForm: form.PostForm{
			Subject: post.Subject,
			Content: post.Content,
			TopicID: post.TopicID,
		},
	}

	v := NewView(w, r, "edit_topic")
	v.Set("form", topicForm)
	v.Render()
}
