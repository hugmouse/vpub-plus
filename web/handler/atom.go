package handler

//
//import (
//	"encoding/xml"
//	"github.com/gorilla/mux"
//	"net/http"
//	"path"
//	"strconv"
//	"time"
//	"vpub/syntax"
//	"vpub/model"
//)
//
//type TimeStr string
//
//type Feed struct {
//	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
//	Subject    string   `xml:"title"`
//	ID       string   `xml:"id"`
//	Link     []Link   `xml:"link"`
//	Updated  TimeStr  `xml:"updated"`
//	Author   *Person  `xml:"author"`
//	Icon     string   `xml:"icon,omitempty"`
//	Logo     string   `xml:"logo,omitempty"`
//	Subtitle string   `xml:"subtitle,omitempty"`
//	Entry    []*Entry `xml:"entry"`
//}
//
//type Entry struct {
//	Subject     string  `xml:"title"`
//	ID        string  `xml:"id"`
//	Link      []Link  `xml:"link"`
//	Published TimeStr `xml:"published"`
//	Updated   TimeStr `xml:"updated"`
//	Author    *Person `xml:"author"`
//	Summary   *Text   `xml:"summary"`
//	Content   *Text   `xml:"content"`
//}
//
//type Link struct {
//	Rel      string `xml:"rel,attr,omitempty"`
//	Href     string `xml:"href,attr"`
//	Type     string `xml:"type,attr,omitempty"`
//	HrefLang string `xml:"hreflang,attr,omitempty"`
//	Subject    string `xml:"title,attr,omitempty"`
//	Length   uint   `xml:"length,attr,omitempty"`
//}
//
//type Person struct {
//	Name     string `xml:"name"`
//	URI      string `xml:"uri,omitempty"`
//	Email    string `xml:"email,omitempty"`
//	InnerXML string `xml:",innerxml"`
//}
//
//type Text struct {
//	Type string `xml:"type,attr"`
//	Body string `xml:",chardata"`
//}
//
//func Time(t time.Time) TimeStr {
//	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
//}
//
//func createAtomEntryFromPost(post model.Post, u string) *Entry {
//	postLink := u + "/posts/" + strconv.FormatInt(post.Id, 10)
//	return &Entry{
//		Subject: post.Subject,
//		ID:    postLink,
//		Link: []Link{
//			{
//				Rel:  "alternate",
//				Href: postLink,
//				Type: "text/html",
//			},
//		},
//		Updated:   Time(post.CreatedAt),
//		Published: Time(post.CreatedAt),
//		Author: &Person{
//			Name: post.User,
//			URI:  u + "/~" + post.User,
//		},
//		Content: &Text{
//			Type: "html",
//			Body: syntax.Convert(post.Content),
//		},
//	}
//}
//
//func (h *Handler) showFeedView(w http.ResponseWriter, r *http.Request) {
//	u := path.Clean(h.url)
//	feed := Feed{
//		Subject:   h.title,
//		ID:      h.url,
//		Updated: Time(time.Now()),
//		Link: []Link{
//			{
//				Rel:  "self",
//				Href: u + "/feed.atom",
//			},
//			{
//				Rel:  "alternate",
//				Type: "text/html",
//				Href: u,
//			},
//		},
//	}
//	posts, err := h.storage.PostsFeed()
//	if err != nil {
//		serverError(w, err)
//		return
//	}
//	for _, post := range posts {
//		feed.Entry = append(feed.Entry, createAtomEntryFromPost(post, u))
//	}
//	var data []byte
//	data, err = xml.MarshalIndent(&feed, "", "    ")
//	if err != nil {
//		serverError(w, err)
//	}
//	w.Write([]byte(xml.Header + string(data)))
//}
//
//func (h *Handler) showFeedViewTopic(w http.ResponseWriter, r *http.Request) {
//	topic := mux.Vars(r)["topic"]
//	if !contains(h.topics, topic) {
//		notFound(w)
//		return
//	}
//	u := path.Clean(h.url)
//	feed := Feed{
//		Subject:   h.title + " - " + topic,
//		ID:      u + "/topics/" + topic,
//		Updated: Time(time.Now()),
//		Link: []Link{
//			{
//				Rel:  "self",
//				Href: u + "/topics/" + topic + "/feed.atom",
//			},
//			{
//				Rel:  "alternate",
//				Type: "text/html",
//				Href: u + "/topics/" + topic,
//			},
//		},
//	}
//	posts, err := h.storage.PostsTopicFeed(topic)
//	if err != nil {
//		serverError(w, err)
//		return
//	}
//	for _, post := range posts {
//		feed.Entry = append(feed.Entry, createAtomEntryFromPost(post, u))
//	}
//	var data []byte
//	data, err = xml.MarshalIndent(&feed, "", "    ")
//	if err != nil {
//		serverError(w, err)
//	}
//	w.Write([]byte(xml.Header + string(data)))
//}
//
////func (h *Handler) showNewPostViewTopic(w http.ResponseWriter, r *http.Request, user string) {
////	var topic string
////	if val, ok := r.URL.Query()["topic"]; ok && len(val) == 1 {
////		topic = val[0]
////	}
////	if !contains(h.topics, topic) && topic != "" {
////		notFound(w)
////		return
////	}
////	postForm := form.PostForm{}
////	postForm.Topics = h.topics
////	postForm.Topic = topic
////	h.renderLayout(w, "create_post", map[string]interface{}{
////		"form":           postForm,
////		csrf.TemplateTag: csrf.TemplateField(r),
////	}, user)
////}
