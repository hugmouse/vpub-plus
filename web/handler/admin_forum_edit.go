package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminEditForumView(w http.ResponseWriter, r *http.Request) {
	forum, err := h.storage.ForumByID(RouteInt64Param(r, "forumId"))
	if err != nil {
		serverError(w, err)
		return
	}

	groups, err := h.storage.Groups()
	if err != nil {
		serverError(w, err)
		return
	}

	forumForm := form.ForumForm{
		Name:                 forum.Name,
		Position:             forum.Position,
		IsLocked:             forum.IsLocked,
		GroupID:              forum.GroupID,
		RestrictedVisibility: forum.RestrictedVisibility,
		Groups:               groups,
	}

	v := NewView(w, r, "admin_forum_edit")
	v.Set("forum", forum)
	v.Set("form", forumForm)
	v.Render()
}
