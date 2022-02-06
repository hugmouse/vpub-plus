package handler

import (
	"net/http"
)

func (h *Handler) showTopicView(w http.ResponseWriter, r *http.Request) {
	topic, err := h.storage.TopicById(RouteInt64Param(r, "topicId"))
	if err != nil {
		notFound(w)
		return
	}
	board, err := h.storage.BoardById(topic.BoardId)
	posts, _, err := h.storage.PostsByTopicId(topic.Id)

	v := NewView(w, r, "topic")
	v.Set("board", board)
	v.Set("topic", topic)
	v.Set("posts", posts)
	v.Render()
}
