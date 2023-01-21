package handler

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"time"
	"vpub/web/handler/request"
)

func (h *Handler) showTopicFeed(w http.ResponseWriter, r *http.Request) {
	topicId := RouteInt64Param(r, "topicId")
	settings := request.GetSettingsContextKey(r)
	feed := Feed{
		Title:   settings.Name,
		ID:      settings.URL,
		Updated: Time(time.Now()),
		Link: []Link{
			{
				Rel:  "self",
				Href: joinPath(settings.URL, "topics", strconv.FormatInt(topicId, 10), "feed.atom"),
			},
			{
				Rel:  "alternate",
				Type: "text/html",
				Href: joinPath(settings.URL, "topics", strconv.FormatInt(topicId, 10)),
			},
		},
	}

	posts, _, err := h.storage.PostsByTopicId(topicId)
	if err != nil {
		serverError(w, err)
		return
	}

	for _, post := range posts {
		feed.Entry = append(feed.Entry, createAtomEntryFromPost(settings.URL, post))
	}

	w.Header().Set("Content-Type", "application/atom+xml")

	var data []byte
	data, err = xml.MarshalIndent(&feed, "", "    ")
	if err != nil {
		serverError(w, err)
	}

	_, err = w.Write([]byte(xml.Header + string(data)))
	if err != nil {
		serverError(w, err)
	}
}
