package handler

import (
	"net/http"
	"vpub/web/handler/request"
)

func (h *Handler) removeAdminForum(w http.ResponseWriter, r *http.Request) {
	forum, err := h.storage.ForumById(RouteInt64Param(r, "forumId"))
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "admin_forum_remove")
	v.Set("forum", forum)

	if err := h.storage.RemoveForum(forum.Id); err != nil {
		v.Set("errorMessage", "Unable to delete forum: "+err.Error())
		v.Render()
		return
	}

	sess := request.GetSessionContextKey(r)
	sess.FlashInfo("Successfully deleted forum")
	sess.Save(r, w)

	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}
