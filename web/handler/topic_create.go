package handler

import (
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) showCreateTopicView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardByID(id)
	if err != nil {
		notFound(w)
		return
	}

	user := request.GetUserContextKey(r)
	if !canAccessForum(board.Forum, user) {
		forbidden(w)
		return
	}

	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	var visibleBoards []model.Board
	for _, b := range boards {
		if canAccessForum(b.Forum, user) {
			visibleBoards = append(visibleBoards, b)
		}
	}
	topicForm := form.TopicForm{
		BoardID: board.ID,
		Boards:  visibleBoards,
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
