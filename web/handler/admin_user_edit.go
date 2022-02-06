package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminEditUserView(w http.ResponseWriter, r *http.Request) {
	user, err := h.storage.UserById(RouteInt64Param(r, "userId"))
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "admin_user_edit")
	v.Set("user", user)
	v.Set("form", form.AdminUserForm{
		Username: user.Name,
		About:    user.About,
		Picture:  user.Picture,
	})
	v.Render()
}
