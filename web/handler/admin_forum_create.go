package handler

import (
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminCreateForumView(w http.ResponseWriter, r *http.Request) {
	groups, err := h.storage.Groups()
	if err != nil {
		serverError(w, err)
		return
	}
	forumForm := form.ForumForm{
		RestrictedVisibility: model.RestrictedVisibilityHidden,
		Groups:               groups,
	}
	v := NewView(w, r, "admin_forum_create")
	v.Set("form", forumForm)
	v.Render()
}
