package handler

import "net/http"

func (h *Handler) showAdminGroupListView(w http.ResponseWriter, r *http.Request) {
	groups, err := h.storage.Groups()
	if err != nil {
		serverError(w, err)
		return
	}
	v := NewView(w, r, "admin_group")
	v.Set("groups", groups)
	v.Render()
}
