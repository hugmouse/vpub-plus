package handler

import (
	"net/http"
	"strconv"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 0
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	user := h.session.Get(r)

	posts, _, err := h.storage.Posts(page, 50)
	if err != nil {
		serverError(w, err)
		return
	}

	users, err := h.storage.RandomUsers(20)
	if err != nil {
		serverError(w, err)
		return
	}

	h.renderLayout(w, "index", map[string]interface{}{
		"posts":            posts,
		"hasNotifications": user.HasNotification,
		"name":             user.Name,
		"users":            users,
	}, "") // TODO
}
