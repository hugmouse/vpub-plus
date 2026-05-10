package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"vpub/web/handler/request"
)

func (h *Handler) showAdminRemoveGroupView(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "groupId")
	group, err := h.storage.GroupByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFound(w)
			return
		}
		serverError(w, err)
		return
	}
	count, err := h.storage.ForumsWithGroupCount(id)
	if err != nil {
		serverError(w, err)
		return
	}
	v := NewView(w, r, "admin_group_remove")
	v.Set("group", group)
	v.Set("affectedForums", count)
	v.Render()
}

func (h *Handler) removeAdminGroup(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "groupId")
	if err := h.storage.RemoveGroup(id); err != nil {
		serverError(w, err)
		return
	}
	session := request.GetSessionContextKey(r)
	session.FlashInfo("Group deleted")
	session.Save(r, w)
	http.Redirect(w, r, "/admin/groups", http.StatusFound)
}
