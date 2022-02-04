package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
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
