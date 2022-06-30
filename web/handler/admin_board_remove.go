package handler

import (
	"net/http"
	"vpub/web/handler/request"
)

func (h *Handler) removeAdminBoard(w http.ResponseWriter, r *http.Request) {
	board, err := h.storage.BoardById(RouteInt64Param(r, "boardId"))
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "admin_board_remove")
	v.Set("board", board)

	if err := h.storage.RemoveBoard(board.Id); err != nil {
		v.Set("errorMessage", "Unable to delete board: "+err.Error())
		v.Render()
		return
	}

	sess := request.GetSessionContextKey(r)
	sess.FlashInfo("Successfully deleted board")
	sess.Save(r, w)

	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}
