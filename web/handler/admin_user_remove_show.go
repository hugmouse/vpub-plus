package handler

import (
	"net/http"
)

func (h *Handler) showAdminRemoveUserView(w http.ResponseWriter, r *http.Request) {
	user, err := h.storage.UserById(RouteInt64Param(r, "userId"))
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "admin_user_remove")
	v.Set("user", user)
	v.Render()
}
