package handler

import (
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
	if err != nil {
		serverError(w, err)
		return
	}

	postForm := form.PostForm{
		Subject: post.Subject,
		Content: post.Content,
		TopicId: post.TopicId,
	}

	v := NewView(w, r, "edit_post")
	v.Set("navigation", navigation{
		Forum: board.Forum,
		Board: board,
		Topic: topic.Post.Subject,
	})
	v.Set("form", postForm)
	v.Set("post", post)
	v.Set("topic", topic)
	v.Set("board", board)
	v.Render()
}
