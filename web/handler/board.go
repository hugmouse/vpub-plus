package handler

import (
	"github.com/gorilla/csrf"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

func (h *Handler) showNewTopicView(w http.ResponseWriter, r *http.Request, user string) {
	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}
	topicForm := form.TopicForm{}
	h.renderLayout(w, "create_topic", map[string]interface{}{
		"form":           topicForm,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) saveTopic(w http.ResponseWriter, r *http.Request, user string) {
	topicForm := form.NewTopicForm(r)
	topic := model.Topic{
		BoardId: topicForm.BoardId,
		FirstPost: model.Post{
			User:    user,
			Title:   topicForm.PostForm.Subject,
			Content: topicForm.PostForm.Content,
		},
	}
	if err := topic.FirstPost.Validate(); err != nil {
		serverError(w, err)
		return
	}
	_, err := h.storage.CreateTopic(topic)
	if err != nil {
		serverError(w, err)
		return
	}
	//http.Redirect(w, r, fmt.Sprintf("/posts/%d", id), http.StatusFound)
}

func (h *Handler) showBoardView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)

	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}

	threads, _, err := h.storage.ThreadsByTopicId(board.Id)

	h.renderLayout(w, "board", map[string]interface{}{
		"board":   board,
		"threads": threads,
		"hasMore": "",
	}, user)
}

//
//func (h *Handler) showPageNumber(w http.ResponseWriter, r *http.Request) {
//	user, _ := h.session.Get(r)
//
//	var page int64 = 0
//	if val, ok := mux.Vars(r)["nb"]; ok {
//		page, _ = strconv.ParseInt(val, 10, 64)
//	}
//
//	var topic string
//	if val, ok := r.URL.Query()["topic"]; ok && len(val) == 1 {
//		topic = val[0]
//	}
//	if !contains(h.topics, topic) && topic != "" {
//		notFound(w)
//		return
//	}
//
//	var posts []model.Post
//	var hasMore bool
//	var err error
//	if topic != "" {
//		posts, hasMore, err = h.storage.PostsTopicWithReplyCount(topic, page, h.perPage)
//	} else {
//		posts, hasMore, err = h.storage.PostsWithReplyCount(page, h.perPage)
//	}
//
//	if err != nil {
//		serverError(w, err)
//		return
//	}
//
//	h.renderLayout(w, "paginate", map[string]interface{}{
//		"topic":    topic,
//		"posts":    posts,
//		"page":     page,
//		"topics":   h.topics,
//		"hasMore":  hasMore,
//		"nextPage": page + 1,
//	}, user)
//}
