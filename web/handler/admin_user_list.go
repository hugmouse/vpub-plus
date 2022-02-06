package handler

import "net/http"

func (h *Handler) showAdminUserListView(w http.ResponseWriter, r *http.Request) {
	users, err := h.storage.Users()
	if err != nil {
		serverError(w, err)
		return
	}

	v := NewView(w, r, "admin_user")
	v.Set("users", users)
	v.Render()
}
