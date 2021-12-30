package handler

import "net/http"

func (h *Handler) showStylesheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(h.css)
}
