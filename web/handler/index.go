package handler

import (
	"html/template"
	"net/http"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "index", map[string]interface{}{
		"posts":            "",
		"hasNotifications": false,
		"boards":           boards,
		"motd":             template.HTML(h.motd),
		"logged":           user,
		"boardTitle":       h.title,
	}, user)
}
