package handler

import "net/http"

func (h *Handler) showAdminForumsView(w http.ResponseWriter, r *http.Request) {
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, r, "admin_forum", map[string]interface{}{
		"forums": forums,
	})
}
