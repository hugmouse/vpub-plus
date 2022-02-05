package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"vpub/config"
	"vpub/model"
	"vpub/storage"
	"vpub/web/handler/request"
	"vpub/web/session"
)

func RouteInt64Param(r *http.Request, param string) int64 {
	vars := mux.Vars(r)
	value, err := strconv.ParseInt(vars[param], 10, 64)
	if err != nil {
		return 0
	}

	if value < 0 {
		return 0
	}

	return value
}

func forbidden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func notFound(w http.ResponseWriter) {
	http.Error(w, "Page Not Found", http.StatusNotFound)
}

func serverError(w http.ResponseWriter, err error) {
	log.Println("[server error]", err)
	http.Error(w, fmt.Sprintf("server error: %s", err), http.StatusInternalServerError)
}

func (h *Handler) handleSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := h.session.GetUser(r)
		session, err := h.session.GetSession(r)
		if err != nil {
			fmt.Println("Unable to create session")
		}
		settings, err := h.storage.Settings()
		if err != nil {
			serverError(w, err)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, request.SessionKey, session)
		ctx = context.WithValue(ctx, request.UserKey, user)
		ctx = context.WithValue(ctx, request.SettingsKey, settings)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type Handler struct {
	session *session.Manager
	url     string
	env     string
	mux     *mux.Router
	storage *storage.Storage
	perPage int
}

func (h *Handler) protect(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := request.GetUserContextKey(r)
		if user.Name == "" {
			forbidden(w)
			return
		}
		fn(w, r)
	}
}

func (h *Handler) admin(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := request.GetUserContextKey(r)
		if !user.IsAdmin {
			forbidden(w)
			return
		}
		fn(w, r)
	}
}

type pagination struct {
	HasMore bool
	Page    int64
}

func forumFromBoards(boards []model.Board) []model.Forum {
	var forums []model.Forum
	var forum model.Forum
	for i, board := range boards {
		if i == 0 {
			forum.Name = board.Forum.Name
			forum.Id = board.Forum.Id
		} else if board.Forum.Id != forum.Id {
			forums = append(forums, forum)
			forum = model.Forum{Name: board.Forum.Name, Id: board.Forum.Id}
		}
		forum.Boards = append(forum.Boards, board)
	}
	if len(forum.Boards) > 0 {
		forums = append(forums, forum)
	}
	return forums
}

func New(cfg *config.Config, data *storage.Storage, s *session.Manager) (http.Handler, error) {
	router := mux.NewRouter()
	h := &Handler{
		session: s,
		mux:     router,
		storage: data,
	}
	router.Use(h.handleSessionMiddleware)
	h.initTpl()

	// Static assets
	router.HandleFunc("/style.css", h.showStylesheet).Methods(http.MethodGet)

	// All public views
	publicSubRouter := router.PathPrefix("/").Subrouter()

	//router.HandleFunc("/favicon.ico", h.showFavicon).Name("favicon").Methods(http.MethodGet)
	//router.HandleFunc("/feed.atom", h.showFeedView).Methods(http.MethodGet)

	// Auth
	publicSubRouter.HandleFunc("/login", h.showLoginView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/login", h.checkLogin).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/register", h.showRegisterView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/register", h.register).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/logout", h.logout).Methods(http.MethodGet)

	publicSubRouter.HandleFunc("/feed.atom", h.showFeed).Methods(http.MethodGet)

	// Boards
	publicSubRouter.HandleFunc("/boards/{boardId}", h.showBoardView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/boards/{boardId}/new-topic", h.protect(h.showCreateTopicView)).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/boards/{boardId}/save-topic", h.protect(h.saveTopic)).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/boards/{boardId}/newest", h.showNewestBoardView).Methods(http.MethodGet)

	// Forums
	publicSubRouter.HandleFunc("/forums/{forumId}", h.showForumView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/posts", h.showPostListView).Methods(http.MethodGet)

	// Topic
	publicSubRouter.HandleFunc("/topics/{topicId}", h.showTopicView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/topics/{topicId}/edit", h.protect(h.showEditTopicView)).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/topics/{topicId}/update", h.admin(h.updateTopic)).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/topics/{topicId}/newest", h.showNewestTopicView).Methods(http.MethodGet)

	publicSubRouter.HandleFunc("/posts/save", h.protect(h.savePost)).Methods(http.MethodPost)
	//publicSubRouter.HandleFunc("/posts/{postId}", h.showPostView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/posts/{postId}/edit", h.protect(h.showEditPostView)).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/posts/{postId}/update", h.protect(h.updatePost)).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/posts/{postId}/remove", h.protect(h.removePost))
	//publicSubRouter.HandleFunc("/posts/{postId}/reply", h.protect(h.savePostReply)).Methods(http.MethodPost)

	// Pagination
	//publicSubRouter.HandleFunc("/page/{nb}", h.showPageNumber).Methods(http.MethodGet)

	// Posts
	//publicSubRouter.HandleFunc("/replies/{replyId}", h.protect(h.showReplyView)).Methods(http.MethodGet)
	//publicSubRouter.HandleFunc("/replies/{replyId}/save", h.protect(h.saveReplyReply)).Methods(http.MethodPost)
	//publicSubRouter.HandleFunc("/replies/{replyId}/edit", h.protect(h.showEditReplyView)).Methods(http.MethodGet)
	//publicSubRouter.HandleFunc("/replies/{replyId}/update", h.protect(h.updateReply)).Methods(http.MethodPost)
	//publicSubRouter.HandleFunc("/replies/{replyId}/remove", h.protect(h.handleRemoveReply))

	// Notifications
	//publicSubRouter.HandleFunc("/notifications", h.protect(h.showNotificationsView)).Methods(http.MethodGet)
	//publicSubRouter.HandleFunc("/notifications/{notificationId}/mark-read", h.protect(h.markRead)).Methods(http.MethodGet)
	//publicSubRouter.HandleFunc("/notifications/mark-all-read", h.protect(h.markAllRead)).Methods(http.MethodGet)

	// User
	//publicSubRouter.HandleFunc("/~{userId}", h.showUserPostsView).Methods(http.MethodGet)
	//publicSubRouter.HandleFunc("/account", h.protect(h.showAccountView)).Methods(http.MethodGet)
	//publicSubRouter.HandleFunc("/save-account", h.protect(h.saveAccount)).Methods(http.MethodPost)

	publicSubRouter.HandleFunc("/reset-password", h.showResetPasswordView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/reset-password", h.updatePassword).Methods(http.MethodPost)

	// Index
	publicSubRouter.HandleFunc("/", h.showIndexView).Name("index").Methods(http.MethodGet)

	adminSubRouter := router.PathPrefix("/admin").Subrouter().StrictSlash(true)
	// Admin
	adminSubRouter.HandleFunc("/", h.admin(h.showAdminView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards", h.admin(h.showAdminBoardsView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards/new", h.admin(h.showAdminCreateBoardView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards/save", h.admin(h.saveAdminBoard)).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/boards/{boardId}/edit", h.admin(h.showAdminEditBoardView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards/{boardId}/update", h.admin(h.updateAdminBoard)).Methods(http.MethodPost)

	adminSubRouter.HandleFunc("/users", h.admin(h.showAdminUserListView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/users/{name}/edit", h.admin(h.showAdminEditUserView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/users/{name}/update", h.admin(h.updateAdminUser)).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/users/{userId}/remove", h.admin(h.showAdminRemoveUserView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/users/{userId}/remove", h.admin(h.removeAdminUser)).Methods(http.MethodPost)

	adminSubRouter.HandleFunc("/settings/edit", h.admin(h.showAdminSettingsView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/settings/update", h.admin(h.updateAdminSettings)).Methods(http.MethodPost)

	adminSubRouter.HandleFunc("/keys", h.admin(h.showAdminKeyListView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/keys/save", h.admin(h.saveAdminKey)).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/keys/{keyId}/remove", h.admin(h.removeAdminKey))

	adminSubRouter.HandleFunc("/forums", h.admin(h.showAdminForumsView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/forums/new", h.admin(h.showAdminCreateForumView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/forums/save", h.admin(h.saveAdminForum)).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/forums/{forumId}/edit", h.admin(h.showAdminEditForumView)).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/forums/{forumId}/update", h.admin(h.updateAdminForum)).Methods(http.MethodPost)

	return router, nil
}
