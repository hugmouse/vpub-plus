package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) showResetPasswordView(w http.ResponseWriter, r *http.Request) {
	var hash string
	if val, ok := r.URL.Query()["hash"]; ok && len(val) == 1 {
		hash = val[0]
	}
	if !h.storage.UserHashExists(hash) {
		notFound(w)
		return
	}
	h.renderLayout(w, r, "reset_password", map[string]interface{}{
		"hash":           hash,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

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
