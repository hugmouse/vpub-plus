package handler

import (
	"fmt"
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

func (h *Handler) showNewTopicView(w http.ResponseWriter, r *http.Request, user model.User) {
	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}
	boards, err := h.storage.Boards()
	topicForm := form.TopicForm{
		BoardId: board.Id,
		Boards:  boards,
	}
	h.renderLayout(w, "create_topic", map[string]interface{}{
		"form":           topicForm,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) saveTopic(w http.ResponseWriter, r *http.Request, user model.User) {
	topicForm := form.NewTopicForm(r)
	post := model.Post{
		User:    user,
		Subject: topicForm.PostForm.Subject,
		Content: topicForm.PostForm.Content,
	}
	if err := post.Validate(); err != nil {
		serverError(w, err)
		return
	}
	boardId := topicForm.NewBoardId
	if boardId == 0 {
		boardId = topicForm.BoardId
	}
	id, err := h.storage.CreateTopic(model.Topic{
		BoardId:  boardId,
		IsSticky: topicForm.IsSticky,
		IsLocked: topicForm.IsLocked,
		Post:     post,
	})
	if err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), http.StatusFound)
}

func (h *Handler) showForumView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)

	id := RouteInt64Param(r, "forumId")
	forum, err := h.storage.ForumById(id)
	if err != nil {
		notFound(w)
		return
	}

	boards, err := h.storage.BoardsByForumId(forum.Id)
	if err != nil {
		notFound(w)
		return
	}

	h.renderLayout(w, "boards", map[string]interface{}{
		"forum":  forum,
		"boards": boards,
	}, user)
}

func (h *Handler) showPostsView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)

	posts, _, err := h.storage.Posts()
	if err != nil {
		notFound(w)
		return
	}

	h.renderLayout(w, "posts", map[string]interface{}{
		"posts":   posts,
		"hasMore": "",
	}, user)
}

func (h *Handler) showBoardView(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)

	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}

	topics, _, err := h.storage.TopicsByBoardId(board.Id)
	if err != nil {
		notFound(w)
		return
	}

	h.renderLayout(w, "board", map[string]interface{}{
		"board":   board,
		"topics":  topics,
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
