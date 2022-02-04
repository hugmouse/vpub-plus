package handler

import (
	"net/http"
	"strconv"
)

func (h *Handler) showPostListView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	posts, hasMore, err := h.storage.Posts(page)
	if err != nil {
		notFound(w)
		return
	}

	h.renderLayout(w, r, "posts", map[string]interface{}{
		"posts": posts,
		"pagination": pagination{
			HasMore: hasMore,
			Page:    page,
		},
	})
}
