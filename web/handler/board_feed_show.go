package handler

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"vpub/model"
	"vpub/syntax"
	"vpub/web/handler/request"
)

func createAtomEntryFromTopic(url string, topic model.Topic, renderEngine syntax.Renderer) *Entry {
	link := joinPath(
		url,
		"topics",
	) + fmt.Sprintf("/%d", topic.Id)
	postCount := ""
	if topic.Posts > 2 {
		postCount = fmt.Sprintf("<p>%d replies</p>", topic.Posts-1)
	} else if topic.Posts == 2 {
		postCount = "<p>1 reply</p>"
	}
	return &Entry{
		Title: topic.Post.Subject,
		ID:    link,
		Link: []Link{
			{
				Rel:  "alternate",
				Href: link,
				Type: "text/html",
			},
		},
		Updated:   Time(topic.UpdatedAt),
		Published: Time(topic.Post.CreatedAt),
		Author: &Person{
			Name: topic.Post.User.Name,
		},
		Content: &Text{
			Type: "html",
			Body: renderEngine.Convert(topic.Post.Content, true) + postCount,
		},
	}
}

func (h *Handler) showBoardFeed(w http.ResponseWriter, r *http.Request) {
	boardId := RouteInt64Param(r, "boardId")
	settings := request.GetSettingsContextKey(r)
	feed := Feed{
		Title:   settings.Name,
		ID:      settings.URL,
		Updated: Time(time.Now()),
		Link: []Link{
			{
				Rel:  "self",
				Href: joinPath(settings.URL, "boards", strconv.FormatInt(boardId, 10), "feed.atom"),
			},
			{
				Rel:  "alternate",
				Type: "text/html",
				Href: joinPath(settings.URL, "boards", strconv.FormatInt(boardId, 10)),
			},
		},
	}

	topics, _, err := h.storage.TopicsByBoardId(boardId, 1)
	if err != nil {
		serverError(w, err)
		return
	}

	for _, topic := range topics {
		feed.Entry = append(feed.Entry, createAtomEntryFromTopic(settings.URL, topic, *h.currentRenderEngine))
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
