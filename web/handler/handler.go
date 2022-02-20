package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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

type navigation struct {
	Forum model.Forum
	Board model.Board
	Topic string
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

func (h *Handler) handleAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := request.GetUserContextKey(r)
		if !user.IsAdmin {
			forbidden(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func New(data *storage.Storage, s *session.Manager) (http.Handler, error) {
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
	//router.HandleFunc("/favicon.ico", h.showFavicon).Name("favicon").Methods(http.MethodGet)

	// Forum views
	publicSubRouter := router.PathPrefix("/").Subrouter()

	// Auth
	publicSubRouter.HandleFunc("/login", h.showLoginView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/login", h.checkLogin).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/register", h.showRegisterView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/register", h.register).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/logout", h.logout).Methods(http.MethodGet)

	// Feed
	publicSubRouter.HandleFunc("/feed.atom", h.showFeed).Methods(http.MethodGet)

	// Forums
	publicSubRouter.HandleFunc("/forums/{forumId}", h.showForumView).Methods(http.MethodGet)

	// Boards
	publicSubRouter.HandleFunc("/boards/{boardId}", h.showBoardView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/boards/{boardId}/feed.atom", h.showBoardFeed).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/boards/{boardId}/new-topic", h.protect(h.showCreateTopicView)).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/boards/{boardId}/save-topic", h.protect(h.saveTopic)).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/boards/{boardId}/newest", h.showNewestBoardView).Methods(http.MethodGet)

	// Topic
	publicSubRouter.HandleFunc("/topics/{topicId}", h.showTopicView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/topics/{topicId}/feed.atom", h.showTopicFeed).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/topics/{topicId}/edit", h.protect(h.showEditTopicView)).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/topics/{topicId}/update", h.updateTopic).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/topics/{topicId}/newest", h.showNewestTopicView).Methods(http.MethodGet)

	// Post
	publicSubRouter.HandleFunc("/posts/save", h.protect(h.savePost)).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/posts/{postId}/edit", h.protect(h.showEditPostView)).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/posts/{postId}/update", h.protect(h.updatePost)).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/posts/{postId}/remove", h.protect(h.removePost))
	publicSubRouter.HandleFunc("/posts", h.showPostListView).Methods(http.MethodGet)

	// Account
	publicSubRouter.HandleFunc("/account", h.protect(h.showAccountEditPage)).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/update-account", h.protect(h.updateAccount)).Methods(http.MethodPost)
	publicSubRouter.HandleFunc("/reset-password", h.showResetPasswordView).Methods(http.MethodGet)
	publicSubRouter.HandleFunc("/reset-password", h.updatePassword).Methods(http.MethodPost)

	// Users
	publicSubRouter.HandleFunc("/users/{userId}", h.showUserView).Methods(http.MethodGet)

	// Index
	publicSubRouter.HandleFunc("/", h.showIndexView).Name("index").Methods(http.MethodGet)

	// Admin router
	adminSubRouter := router.PathPrefix("/admin").Subrouter().StrictSlash(true)
	adminSubRouter.Use(h.handleAdminMiddleware)

	// Admin
	adminSubRouter.HandleFunc("/", h.showAdminView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards", h.showAdminBoardsView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards/new", h.showAdminCreateBoardView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards/save", h.saveAdminBoard).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/boards/{boardId}/edit", h.showAdminEditBoardView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards/{boardId}/update", h.updateAdminBoard).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/boards/{boardId}/remove", h.showAdminRemoveBoardView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/boards/{boardId}/remove", h.removeAdminBoard).Methods(http.MethodPost)

	adminSubRouter.HandleFunc("/users", h.showAdminUserListView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/users/{userId}/edit", h.showAdminEditUserView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/users/{userId}/update", h.updateAdminUser).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/users/{userId}/remove", h.showAdminRemoveUserView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/users/{userId}/remove", h.removeAdminUser).Methods(http.MethodPost)

	adminSubRouter.HandleFunc("/settings/edit", h.showAdminSettingsView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/settings/update", h.updateAdminSettings).Methods(http.MethodPost)

	adminSubRouter.HandleFunc("/keys", h.showAdminKeyListView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/keys/save", h.saveAdminKey).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/keys/{keyId}/remove", h.removeAdminKey)

	adminSubRouter.HandleFunc("/forums", h.showAdminForumsView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/forums/new", h.showAdminCreateForumView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/forums/save", h.saveAdminForum).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/forums/{forumId}/edit", h.showAdminEditForumView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/forums/{forumId}/update", h.updateAdminForum).Methods(http.MethodPost)
	adminSubRouter.HandleFunc("/forums/{forumId}/remove", h.showAdminRemoveForumView).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/forums/{forumId}/remove", h.removeAdminForum).Methods(http.MethodPost)

	return router, nil
}
