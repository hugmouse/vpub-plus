package handler

import (
	"github.com/gorilla/csrf"
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
	topic, err := h.storage.TopicById(RouteInt64Param(r, "topicId"))
	if err != nil {
		notFound(w)
		return
	}
	post, err := h.storage.PostById(topic.Post.Id)
	if err != nil {
		notFound(w)
		return
	}
	boards, err := h.storage.Boards()
	topicForm := form.TopicForm{
		Id:       topic.Id,
		BoardId:  topic.BoardId,
		IsSticky: topic.IsSticky,
		IsLocked: topic.IsLocked,
		Boards:   boards,
		PostForm: form.PostForm{
			Subject: post.Subject,
			Content: post.Content,
			TopicId: post.TopicId,
		},
	}
	h.renderLayout(w, r, "edit_topic", map[string]interface{}{
		"form":           topicForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
