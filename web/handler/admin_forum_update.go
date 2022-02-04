package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) updateAdminForum(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "forumId")
	forum, err := h.storage.ForumById(id)
	if err != nil {
		notFound(w)
		return
	}

	forumForm := form.NewForumForm(r)

	v := NewView(w, r, "admin_forum_edit")
	v.Set("forum", forum)
	v.Set("form", forumForm)

	forumRequest := model.ForumRequest{
		Name:     forumForm.Name,
		Position: forumForm.Position,
		IsLocked: forumForm.IsLocked,
	}

	if err := validator.ValidateForumModification(h.storage, id, forumRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if err := h.storage.UpdateForum(id, forumRequest); err != nil {
		v.Set("errorMessage", "Unable to update forum")
		serverError(w, err)
		return
	}

	session := request.GetSessionContextKey(r)
	session.FlashInfo("Forum updated")
	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/admin/forums/%d/edit", id), http.StatusFound)
}
