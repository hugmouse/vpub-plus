package handler

import "net/http"

func (h *Handler) showAdminView(w http.ResponseWriter, r *http.Request) {
	v := NewView(w, r, "admin")
	v.Render()
}
