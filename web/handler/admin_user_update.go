package handler

import (
	"fmt"
	"net/http"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) updateAdminUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.storage.UserById(RouteInt64Param(r, "userId"))
	if err != nil {
		notFound(w)
		return
	}

	userForm := form.NewAdminUserForm(r)
	user.Name = userForm.Username
	user.About = userForm.About
	user.Picture = userForm.Picture

	if err := h.storage.UpdateUser(user); err != nil {
		serverError(w, err)
		return
	}

	sess := request.GetSessionContextKey(r)
	sess.FlashInfo("User account updated")
	sess.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/admin/users/%d/edit", user.Id), http.StatusFound)
}
