package handler

import (
	"net/http"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) updateAccount(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)

	accountForm := form.NewAccountForm(r)

	v := NewView(w, r, "account")
	v.Set("form", accountForm)

	user.About = accountForm.About
	user.Picture = accountForm.Picture
	user.PictureAlt = accountForm.PictureAlt

	if err := h.storage.UpdateUser(user); err != nil {
		serverError(w, err)
		return
	}

	sess := request.GetSessionContextKey(r)
	sess.FlashInfo("Account settings updated")
	sess.Save(r, w)

	v.Render()
}
