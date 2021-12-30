package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handler) showUserPostsView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	user, err := h.storage.UserByName(mux.Vars(r)["userId"])
	if err != nil {
		notFound(w)
		return
	}

	posts, showMore, err := h.storage.PostsByUsername(user.Name, h.perPage, page)
	if err != nil {
		serverError(w, err)
		return
	}

	h.renderLayout(w, "user_posts", map[string]interface{}{
		"user":     user,
		"posts":    posts,
		"page":     page,
		"showMore": showMore,
		"nextPage": page + 1,
	}, "")
}
