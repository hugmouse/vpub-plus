package handler

import (
	"net/http"
)

func (h *Handler) showAdminKeyListView(w http.ResponseWriter, r *http.Request) {
	keys, err := h.storage.Keys()
	if err != nil {
		serverError(w, err)
		return
	}

	v := NewView(w, r, "admin_keys")
	v.Set("keys", keys)
	v.Render()
}
