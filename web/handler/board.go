package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
	"net/http"
	"strconv"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

func (h *Handler) showNewTopicView(w http.ResponseWriter, r *http.Request) {
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
	h.renderLayout(w, r, "create_topic", map[string]interface{}{
		"form":           topicForm,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) saveTopic(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)

	topicForm := form.NewTopicForm(r)

	boards, err := h.storage.Boards()
	if err != nil {
		notFound(w)
		return
	}

	topicForm.Boards = boards

	board, err := h.storage.BoardById(topicForm.BoardId)
	if err != nil {
		notFound(w)
		return
	}

	v := NewView(w, r, "create_topic")
	v.Set("form", topicForm)
	v.Set("board", board)

	boardId := topicForm.NewBoardId
	if boardId == 0 {
		boardId = topicForm.BoardId
	}

	topicRequest := model.TopicRequest{
		BoardId:  boardId,
		IsSticky: topicForm.IsSticky,
		IsLocked: topicForm.IsLocked,
		UserId:   user.Id,
		Subject:  topicForm.PostForm.Subject,
		Content:  topicForm.PostForm.Content,
	}

	if err := validator.ValidateTopicRequest(topicRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	id, err := h.storage.CreateTopic(topicRequest)
	if err != nil {
		v.Set("errorMessage", "Unable to create topic")
		v.Render()
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), http.StatusFound)
}

func (h *Handler) showForumView(w http.ResponseWriter, r *http.Request) {
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

	h.renderLayout(w, r, "boards", map[string]interface{}{
		"forum":  forum,
		"boards": boards,
	})
}

func (h *Handler) showPostsView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	posts, hasMore, err := h.storage.Posts(page)
	if err != nil {
		notFound(w)
		return
	}

	h.renderLayout(w, r, "posts", map[string]interface{}{
		"posts": posts,
		"pagination": pagination{
			HasMore: hasMore,
			Page:    page,
		},
	})
}

func (h *Handler) showBoardView(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	if val, ok := r.URL.Query()["page"]; ok && len(val[0]) == 1 {
		page, _ = strconv.ParseInt(val[0], 10, 64)
	}

	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		notFound(w)
		return
	}

	topics, hasMore, err := h.storage.TopicsByBoardId(board.Id, page)
	if err != nil {
		notFound(w)
		return
	}

	h.renderLayout(w, r, "board", map[string]interface{}{
		"board":  board,
		"topics": topics,
		"pagination": pagination{
			HasMore: hasMore,
			Page:    page,
		},
	})
}

//
//func (h *Handler) showPageNumber(w http.ResponseWriter, r *http.Request) {
//	user, _ := h.session.GetUser(r)
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
