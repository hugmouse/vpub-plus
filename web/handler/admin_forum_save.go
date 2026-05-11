package handler

import (
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
)

func (h *Handler) saveAdminForum(w http.ResponseWriter, r *http.Request) {
	forumForm, err := form.NewForumForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	groups, err := h.storage.Groups()
	if err != nil {
		serverError(w, err)
		return
	}
	forumForm.Groups = groups

	v := NewView(w, r, "admin_forum_create")
	v.Set("form", forumForm)

	forumRequest := model.ForumRequest{
		Name:                 forumForm.Name,
		Position:             forumForm.Position,
		IsLocked:             forumForm.IsLocked,
		GroupID:              forumForm.GroupID,
		RestrictedVisibility: forumForm.RestrictedVisibility,
	}

	if err := validator.ValidateForumCreation(h.storage, forumRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if _, err := h.storage.CreateForum(forumRequest); err != nil {
		serverError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}
