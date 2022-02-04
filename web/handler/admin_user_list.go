package handler

import "net/http"

func (h *Handler) showAdminUserListView(w http.ResponseWriter, r *http.Request) {
	users, err := h.storage.Users()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, r, "admin_user", map[string]interface{}{
		"users": users,
	})
}
