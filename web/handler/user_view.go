package handler

import (
	"net/http"
	"strconv"
)

func (h *Handler) showUserView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	posts, hasMore, err := h.storage.PostsByUserID(RouteInt64Param(r, "userId"), page)
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "posts")
	v.Set("posts", posts)
	v.Set("pagination", pagination{
		HasMore: hasMore,
		Page:    page,
	})
	v.Render()
}
