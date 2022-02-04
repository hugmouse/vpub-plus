package handler

import (
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	pwdForm := form.NewResetPasswordForm(r)
	if err := pwdForm.Validate(); err != nil {
		serverError(w, err)
		return
	}
	if err := h.storage.UpdatePassword(pwdForm.Hash, model.User{Password: pwdForm.Password}); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
