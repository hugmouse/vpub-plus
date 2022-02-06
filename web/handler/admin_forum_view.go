package handler

import "net/http"

func (h *Handler) showAdminForumsView(w http.ResponseWriter, r *http.Request) {
	forums, err := h.storage.Forums()

	if err != nil {
		serverError(w, err)
		return
	}

	v := NewView(w, r, "admin_forum")
	v.Set("forums", forums)
	v.Render()
}
