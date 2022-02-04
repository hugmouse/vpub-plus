package handler

import "net/http"

func (h *Handler) showLoginView(w http.ResponseWriter, r *http.Request) {
	v := NewView(w, r, "login")
	v.Render()
}
