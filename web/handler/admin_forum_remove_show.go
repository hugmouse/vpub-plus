package handler

import "net/http"

func (h *Handler) showAdminRemoveForumView(w http.ResponseWriter, r *http.Request) {
	forum, err := h.storage.ForumById(RouteInt64Param(r, "forumId"))
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "admin_forum_remove")
	v.Set("forum", forum)
	v.Render()
}
