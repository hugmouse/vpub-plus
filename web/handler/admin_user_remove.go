package handler

import (
	"net/http"
	"vpub/web/handler/request"
)

func (h *Handler) removeAdminUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.storage.UserById(RouteInt64Param(r, "userId"))
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "admin_user_remove")
	v.Set("user", user)

	if err := h.storage.RemoveUser(user.Id); err != nil {
		v.Set("errorMessage", "Unable to delete user: "+err.Error())
		v.Render()
		return
	}

	sess := request.GetSessionContextKey(r)
	sess.FlashInfo("Successfully deleted user")
	sess.Save(r, w)

	http.Redirect(w, r, "/admin/users", http.StatusFound)
}
