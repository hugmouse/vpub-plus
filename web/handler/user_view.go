package handler

import (
	"net/http"
	"strconv"
	"vpub/web/handler/request"
)

func (h *Handler) showUserView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	user := request.GetUserContextKey(r)
	posts, hasMore, err := h.storage.PostsByUserID(RouteInt64Param(r, "userId"), page, user.IsAdmin, user.GroupIDs)
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
