package handler

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
	"vpub/model"
	"vpub/syntax"
	"vpub/web/handler/request"
)

type Feed struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title    string   `xml:"title"`
	ID       string   `xml:"id"`
	Link     []Link   `xml:"link"`
	Updated  TimeStr  `xml:"updated"`
	Author   *Person  `xml:"author"`
	Icon     string   `xml:"icon,omitempty"`
	Logo     string   `xml:"logo,omitempty"`
	Subtitle string   `xml:"subtitle,omitempty"`
	Entry    []*Entry `xml:"entry"`
}

type Entry struct {
	Title     string  `xml:"title"`
	ID        string  `xml:"id"`
	Link      []Link  `xml:"link"`
	Published TimeStr `xml:"published"`
	Updated   TimeStr `xml:"updated"`
	Author    *Person `xml:"author"`
	Summary   *Text   `xml:"summary"`
	Content   *Text   `xml:"content"`
}

type Link struct {
	Rel      string `xml:"rel,attr,omitempty"`
	Href     string `xml:"href,attr"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   uint   `xml:"length,attr,omitempty"`
}

type Person struct {
	Name     string `xml:"name"`
	URI      string `xml:"uri,omitempty"`
	Email    string `xml:"email,omitempty"`
	InnerXML string `xml:",innerxml"`
}

type Text struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type TimeStr string

func Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

func joinPath(base string, p ...string) string {
	u, _ := url.Parse(base)
	u.Path = path.Join(p...)
	return u.String()
}

func createAtomEntryFromPost(url string, post model.Post, renderEngine syntax.Renderer) *Entry {
	link := joinPath(
		url,
		"topics",
	) + fmt.Sprintf("/%d#%d", post.TopicId, post.Id)
	return &Entry{
		Title: post.Subject,
		ID:    link,
		Link: []Link{
			{
				Rel:  "alternate",
				Href: link,
				Type: "text/html",
			},
		},
		Updated:   Time(post.UpdatedAt),
		Published: Time(post.CreatedAt),
		Author: &Person{
			Name: post.User.Name,
		},
		Content: &Text{
			Type: "html",
			Body: renderEngine.Convert(post.Content, true),
		},
	}
}

func (h *Handler) showFeed(w http.ResponseWriter, r *http.Request) {
	settings := request.GetSettingsContextKey(r)

	feed := Feed{
		Title:   settings.Name,
		ID:      settings.URL,
		Updated: Time(time.Now()),
		Link: []Link{
			{
				Rel:  "self",
				Href: joinPath(settings.URL, "feed.atom"),
			},
			{
				Rel:  "alternate",
				Type: "text/html",
				Href: settings.URL,
			},
		},
	}

	posts, _, err := h.storage.Posts(1)
	if err != nil {
		serverError(w, err)
		return
	}

	for _, post := range posts {
		feed.Entry = append(feed.Entry, createAtomEntryFromPost(settings.URL, post, *h.currentRenderEngine))
	}

	w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")

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
