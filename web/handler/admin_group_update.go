package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) updateAdminGroup(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "groupId")
	group, err := h.storage.GroupByID(id)
	if err != nil {
		notFound(w)
		return
	}

	groupForm := form.NewGroupForm(r)
	req := model.GroupRequest{Name: groupForm.Name}

	if err := validator.ValidateGroupModification(h.storage, id, req); err != nil {
		members, mErr := h.storage.GroupMembers(id)
		if mErr != nil {
			serverError(w, mErr)
			return
		}
		allUsers, uErr := h.storage.Users()
		if uErr != nil {
			serverError(w, uErr)
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
		v.Set("form", groupForm)
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if err := h.storage.UpdateGroup(id, req); err != nil {
		serverError(w, err)
		return
	}

	session := request.GetSessionContextKey(r)
	session.FlashInfo("Group updated")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/admin/groups/%d/edit", id), http.StatusFound)
}

func (h *Handler) addAdminGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := RouteInt64Param(r, "groupId")
	userID, _ := parseInt64(r.FormValue("user_id"))
	if userID == 0 {
		http.Redirect(w, r, fmt.Sprintf("/admin/groups/%d/edit", groupID), http.StatusFound)
		return
	}
	if err := h.storage.AddGroupMember(groupID, userID); err != nil {
		serverError(w, err)
		return
	}
	session := request.GetSessionContextKey(r)
	session.FlashInfo("Member added")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/admin/groups/%d/edit", groupID), http.StatusFound)
}

func (h *Handler) removeAdminGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := RouteInt64Param(r, "groupId")
	userID := RouteInt64Param(r, "userId")
	if err := h.storage.RemoveGroupMember(groupID, userID); err != nil {
		serverError(w, err)
		return
	}
	session := request.GetSessionContextKey(r)
	session.FlashInfo("Member removed")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/admin/groups/%d/edit", groupID), http.StatusFound)
}

func parseInt64(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
