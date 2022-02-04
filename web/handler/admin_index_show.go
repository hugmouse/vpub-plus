package handler

import "net/http"

func (h *Handler) showAdminView(w http.ResponseWriter, r *http.Request) {
	h.renderLayout(w, r, "admin", nil)
}
