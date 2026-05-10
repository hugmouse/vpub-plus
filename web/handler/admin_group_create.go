package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminCreateGroupView(w http.ResponseWriter, r *http.Request) {
	v := NewView(w, r, "admin_group_create")
	v.Set("form", form.GroupForm{})
	v.Render()
}
