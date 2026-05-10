package handler

import (
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminEditGroupView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "groupId")
	group, err := h.storage.GroupByID(id)
	if err != nil {
		notFound(w)
		return
	}

	members, err := h.storage.GroupMembers(id)
	if err != nil {
		serverError(w, err)
		return
	}

	allUsers, err := h.storage.Users()
	if err != nil {
		serverError(w, err)
		return
	}

	memberSet := make(map[int64]bool, len(members))
	for _, m := range members {
		memberSet[m.ID] = true
	}
	var nonMembers []model.User
	for _, u := range allUsers {
		if !memberSet[u.ID] {
			nonMembers = append(nonMembers, u)
		}
	}

	v := NewView(w, r, "admin_group_edit")
	v.Set("group", group)
	v.Set("members", members)
	v.Set("nonMembers", nonMembers)
	v.Set("form", form.GroupForm{Name: group.Name})
	v.Render()
}
