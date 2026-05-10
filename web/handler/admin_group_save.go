package handler

import (
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
)

func (h *Handler) saveAdminGroup(w http.ResponseWriter, r *http.Request) {
	groupForm := form.NewGroupForm(r)
	v := NewView(w, r, "admin_group_create")
	v.Set("form", groupForm)

	req := model.GroupRequest{Name: groupForm.Name}

	if err := validator.ValidateGroupCreation(h.storage, req); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if _, err := h.storage.CreateGroup(req); err != nil {
		serverError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/groups", http.StatusFound)
}
