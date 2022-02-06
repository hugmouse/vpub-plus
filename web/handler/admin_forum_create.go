package handler

import (
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminCreateForumView(w http.ResponseWriter, r *http.Request) {
	forumForm := form.ForumForm{}
	v := NewView(w, r, "admin_forum_create")
	v.Set("form", forumForm)
	v.Render()
}
