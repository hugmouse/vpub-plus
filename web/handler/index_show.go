package handler

import "net/http"

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	forums := forumFromBoards(boards)
	h.renderLayout(w, r, "index", map[string]interface{}{
		"forums": forums,
	})
}
