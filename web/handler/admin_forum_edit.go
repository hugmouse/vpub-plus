package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminEditForumView(w http.ResponseWriter, r *http.Request) {
	forum, err := h.storage.ForumById(RouteInt64Param(r, "forumId"))
	if err != nil {
		serverError(w, err)
		return
	}
	forumForm := form.ForumForm{
		Name:     forum.Name,
		Position: forum.Position,
		IsLocked: forum.IsLocked,
	}
	h.renderLayout(w, r, "admin_forum_edit", map[string]interface{}{
		"forum":          forum,
		"form":           forumForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
