package handler

import (
	"net/http"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	forums := forumFromBoards(boards)

	v := NewView(w, r, "index")
	v.Set("forums", forums)
	v.Render()
}
