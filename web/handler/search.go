package handler

import (
	"net/http"
)

func (h *Handler) searchShow(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	thing, _ := h.storage.Search(query)
	v := NewView(w, r, "search")
	v.Set("sql", thing)
	v.Set("q", query)
	v.Render()
}
