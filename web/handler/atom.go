package handler

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (h *Handler) showFeedView(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title: h.title,
		Link:  &feeds.Link{Href: "TODO"}, // TODO
		//Description: "Virtual Pub", // TODO
		Created: now,
	}

	posts, _, err := h.storage.Posts(0, 20)
	if err != nil {
		serverError(w, err)
		return
	}

	for _, post := range posts {
		if err != nil {
			serverError(w, err)
			return
		}
		feed.Items = append(feed.Items, &feeds.Item{
			Title:   post.Title,
			Link:    &feeds.Link{Href: fmt.Sprintf("TODO/%d", post.Id)}, // TODO
			Author:  &feeds.Author{Name: post.User},
			Created: post.CreatedAt,
		})
	}
	f, err := feed.ToAtom()
	if err != nil {
		serverError(w, err)
		return
	}
	w.Write([]byte(f))
}

func (h *Handler) showFeedViewTopic(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title: h.title,
		Link:  &feeds.Link{Href: "TODO"}, // TODO
		//Description: "Virtual Pub", // TODO
		Created: now,
	}

	topic := mux.Vars(r)["topic"]
	if !contains(h.topics, topic) {
		notFound(w)
		return
	}

	posts, _, err := h.storage.PostsTopic(topic, 0, 20)
	if err != nil {
		serverError(w, err)
		return
	}

	for _, post := range posts {
		if err != nil {
			serverError(w, err)
			return
		}
		feed.Items = append(feed.Items, &feeds.Item{
			Title:   post.Title,
			Link:    &feeds.Link{Href: fmt.Sprintf("TODO/%d", post.Id)}, // TODO
			Author:  &feeds.Author{Name: post.User},
			Created: post.CreatedAt,
		})
	}
	atom, err := feed.ToAtom()
	if err != nil {
		serverError(w, err)
		return
	}
	w.Write([]byte(atom))
}
