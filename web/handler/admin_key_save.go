package handler

import "net/http"

func (h *Handler) saveAdminKey(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.CreateKey(); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/keys", http.StatusFound)
}
