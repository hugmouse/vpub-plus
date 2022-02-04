package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
)

func (h *Handler) showAdminKeyListView(w http.ResponseWriter, r *http.Request) {
	keys, err := h.storage.Keys()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, r, "admin_keys", map[string]interface{}{
		"keys":           keys,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
