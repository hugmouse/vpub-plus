package handler

import (
	"net/http"
	"strconv"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)

	var page int64 = 0
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

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

	hasNotifs := h.storage.UserHasNotifications(user)

	h.renderLayout(w, "index", map[string]interface{}{
		"posts":            posts,
		"hasNotifications": hasNotifs,
		"users":            users,
	}, user)
}
