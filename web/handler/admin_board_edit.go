package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminEditBoardView(w http.ResponseWriter, r *http.Request) {
	board, err := h.storage.BoardByID(RouteInt64Param(r, "boardId"))
	if err != nil {
		serverError(w, err)
		return
	}
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.BoardForm{
		Name:        board.Name,
		Description: board.Description,
		Position:    board.Position,
		ForumID:     board.Forum.ID,
		Forums:      forums,
		IsLocked:    board.IsLocked,
	}
	v := NewView(w, r, "admin_board_edit")
	v.Set("form", boardForm)
	v.Set("board", board)
	v.Render()
}
