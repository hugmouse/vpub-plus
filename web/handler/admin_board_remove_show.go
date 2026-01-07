package handler

import "net/http"

func (h *Handler) showAdminRemoveBoardView(w http.ResponseWriter, r *http.Request) {
	board, err := h.storage.BoardByID(RouteInt64Param(r, "boardId"))
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "admin_board_remove")
	v.Set("board", board)
	v.Render()
}
