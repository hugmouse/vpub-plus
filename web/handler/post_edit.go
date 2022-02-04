package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showEditPostView(w http.ResponseWriter, r *http.Request) {
	post, err := h.storage.PostById(RouteInt64Param(r, "postId"))
	if err != nil {
		serverError(w, err)
		return
	}

	topic, err := h.storage.TopicById(post.TopicId)
	if err != nil {
		notFound(w)
		return
	}

	board, err := h.storage.BoardById(topic.BoardId)

	postForm := form.PostForm{
		Subject: post.Subject,
		Content: post.Content,
		TopicId: post.TopicId,
	}
	h.renderLayout(w, r, "edit_post", map[string]interface{}{
		"form":           postForm,
		"post":           post,
		"topic":          topic,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
