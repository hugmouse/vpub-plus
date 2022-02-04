package handler

import (
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
)

func (h *Handler) saveAdminForum(w http.ResponseWriter, r *http.Request) {
	forumForm := form.NewForumForm(r)

	v := NewView(w, r, "admin_forum_create")
	v.Set("form", forumForm)

	forumRequest := model.ForumRequest{
		Name:     forumForm.Name,
		Position: forumForm.Position,
		IsLocked: forumForm.IsLocked,
	}

	if err := validator.ValidateForumCreation(h.storage, forumRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if _, err := h.storage.CreateForum(forumRequest); err != nil {
		v.Set("errorMessage", "Unable to create forum")
		serverError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}
