package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) removeAdminKey(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "keyId")

	h.storage.DeleteKey(id)

	http.Redirect(w, r, fmt.Sprintf("/admin/keys"), http.StatusFound)
}
