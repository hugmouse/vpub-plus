package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) removeAdminKey(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "keyId")

	err := h.storage.DeleteKey(id)
	if err != nil {
		serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/keys"), http.StatusFound)
}
