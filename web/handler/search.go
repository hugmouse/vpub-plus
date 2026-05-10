package handler

import (
	"net/http"
	"vpub/model"
	"vpub/web/handler/request"
)

func (h *Handler) searchShow(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results, err := h.storage.Search(query)
	if err != nil {
		serverError(w, err)
		return
	}

	user := request.GetUserContextKey(r)
	var filtered []model.Search
	for _, result := range results {
		switch {
		case result.ForumGroupID == nil:
			filtered = append(filtered, result) // user profile — always visible
		case *result.ForumGroupID == 0:
			filtered = append(filtered, result) // public forum content
		case user.IsAdmin:
			filtered = append(filtered, result) // admins see all restricted content
		case isMember(*result.ForumGroupID, user):
			filtered = append(filtered, result) // restricted forum, user is a member
		}
	}

	v := NewView(w, r, "search")
	v.Set("sql", filtered)
	v.Set("q", query)
	v.Render()
}
