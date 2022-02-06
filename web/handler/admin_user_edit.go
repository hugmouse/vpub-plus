package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminEditUserView(w http.ResponseWriter, r *http.Request) {
	u, err := h.storage.UserByName(mux.Vars(r)["name"])
	if err != nil {
		serverError(w, err)
		return
	}

	v := NewView(w, r, "admin_user_edit")
	v.Set("user", u)
	v.Set("form", form.AdminUserForm{
		Username: u.Name,
		About:    u.About,
	})
	v.Render()
}
