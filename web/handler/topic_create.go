package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showCreateTopicView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}
	boards, err := h.storage.Boards()
	topicForm := form.TopicForm{
		BoardId: board.Id,
		Boards:  boards,
	}

	v := NewView(w, r, "create_topic")
	v.Set("navigation", navigation{
		Forum: board.Forum,
		Board: board,
		Topic: "New topic",
	})
	v.Set("form", topicForm)
	v.Set("board", board)
	v.Render()
}
