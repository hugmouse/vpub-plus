package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminCreateForumView(w http.ResponseWriter, r *http.Request) {
	forumForm := form.ForumForm{}
	h.renderLayout(w, r, "admin_forum_create", map[string]interface{}{
		"form":           forumForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
