package handler

import (
	"html/template"
	"net/http"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)
	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}
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
		"title":            settings.Name,
	}, user)
}
