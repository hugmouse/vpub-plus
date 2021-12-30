package handler

import (
	"html/template"
	"net/http"
	"strconv"
)

func (h *Handler) showIndexView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)

	var page int64 = 0
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	posts, hasMore, err := h.storage.Posts(page, h.perPage)
	if err != nil {
		serverError(w, err)
		return
	}

	users, err := h.storage.Users()
	if err != nil {
		serverError(w, err)
		return
	}

	hasNotifs := h.storage.UserHasNotifications(user)

	h.view("index").Execute(w, map[string]interface{}{
		"posts":            posts,
		"hasNotifications": hasNotifs,
		"users":            users,
		"topics":           h.topics,
		"motd":             template.HTML(h.motd),
		"logged":           user,
		"hasMore":          hasMore,
		"boardTitle":       h.title,
	})
}
