package handler

import (
	"net/http"
)

func (h *Handler) searchShow(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	thing, err := h.storage.Search(query)
	if err != nil {
		serverError(w, err)
		return
	}
	v := NewView(w, r, "search")
	v.Set("sql", thing)
	v.Set("q", query)
	v.Render()
}
