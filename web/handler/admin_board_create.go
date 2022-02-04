package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminCreateBoardView(w http.ResponseWriter, r *http.Request) {
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.BoardForm{
		Forums: forums,
	}
	v := NewView(w, r, "admin_board_create")
	v.Set("form", boardForm)
	v.Render()
}
