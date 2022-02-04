package handler

import "net/http"

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if err := h.session.Delete(w, r); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
