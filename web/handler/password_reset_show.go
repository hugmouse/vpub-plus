package handler

import (
	"net/http"
)

func (h *Handler) showResetPasswordView(w http.ResponseWriter, r *http.Request) {
	var hash string
	if val, ok := r.URL.Query()["hash"]; ok && len(val) == 1 {
		hash = val[0]
	}
	userHashExists, err := h.storage.UserHashExists(hash)
	if err != nil {
		serverError(w, err)
		return
	}
	if !userHashExists {
		notFound(w)
		return
	}

	v := NewView(w, r, "reset_password")
	v.Set("hash", hash)
	v.Render()
}
