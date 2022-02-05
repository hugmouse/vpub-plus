package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) showNewestTopicView(w http.ResponseWriter, r *http.Request) {
	topicId := RouteInt64Param(r, "topicId")

	id, err := h.storage.NewestPostFromTopic(topicId)

	if err != nil {
		notFound(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d#%d", topicId, id), http.StatusFound)
}
