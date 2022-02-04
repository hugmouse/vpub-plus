package handler

import "net/http"

func (h *Handler) showAdminBoardsView(w http.ResponseWriter, r *http.Request) {
	v := NewView(w, r, "admin_board")
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	hasForums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	forums := forumFromBoards(boards)
	v.Set("hasForums", hasForums)
	v.Set("forums", forums)
	v.Render()
}
