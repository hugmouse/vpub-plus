package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) showNewestBoardView(w http.ResponseWriter, r *http.Request) {
	boardId := RouteInt64Param(r, "boardId")

	id, err := h.storage.NewestTopicFromBoard(boardId)

	if err != nil {
		notFound(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d/newest", id), http.StatusFound)
}
