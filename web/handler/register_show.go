package handler

import "net/http"

func (h *Handler) showRegisterView(w http.ResponseWriter, r *http.Request) {
	v := NewView(w, r, "register")
	v.Render()
}
