package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"vpub/assets"
	"vpub/config"
	"vpub/model"
	"vpub/storage"
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

type ProtectedFunc func(http.ResponseWriter, *http.Request, model.User)

type Handler struct {
	session *session.Session
	url     string
	env     string
	css     []byte
	mux     *mux.Router
	storage *storage.Storage
	topics  []string
	perPage int
}

func (h *Handler) protect(fn ProtectedFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.session.Get(r)
		if err != nil || user.Name == "" {
			forbidden(w)
			return
		}
		fn(w, r, user)
	}
}

func (h *Handler) admin(fn ProtectedFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.session.Get(r)
		if err != nil || !user.IsAdmin {
			forbidden(w)
			return
		}
		fn(w, r, user)
	}
}

func New(cfg *config.Config, data *storage.Storage, s *session.Session) (http.Handler, error) {
	router := mux.NewRouter()
	h := &Handler{
		session: s,
		mux:     router,
		storage: data,
		perPage: cfg.PerPage,
		url:     cfg.URL,
	}
	h.initTpl()

	// Read and cache css
	cssFile, _ := os.Open(cfg.CSSFile)
	b, _ := io.ReadAll(cssFile)
	h.css = []byte(assets.AssetsMap["style"] + "\n" + string(b))

	// Static assets
	router.HandleFunc("/style.css", h.showStylesheet).Methods(http.MethodGet)
	//router.HandleFunc("/favicon.ico", h.showFavicon).Name("favicon").Methods(http.MethodGet)
	//router.HandleFunc("/feed.atom", h.showFeedView).Methods(http.MethodGet)

	// Auth
	router.HandleFunc("/login", h.showLoginView).Methods(http.MethodGet)
	router.HandleFunc("/login", h.checkLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.showRegisterView).Methods(http.MethodGet)
	router.HandleFunc("/register", h.register).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.logout).Methods(http.MethodGet)

	// Boards
	router.HandleFunc("/boards/{boardId}", h.showBoardView).Methods(http.MethodGet)
	router.HandleFunc("/boards/{boardId}/new-topic", h.protect(h.showNewTopicView)).Methods(http.MethodGet)
	router.HandleFunc("/boards/{boardId}/save-topic", h.protect(h.saveTopic)).Methods(http.MethodPost)
	//router.HandleFunc("/boards/{boardId}/feed.atom", h.showFeedViewTopic).Methods(http.MethodGet)

	// Forums
	router.HandleFunc("/forums/{forumId}", h.showForumView).Methods(http.MethodGet)
	router.HandleFunc("/posts", h.showPostsView).Methods(http.MethodGet)

	// Topic
	router.HandleFunc("/topics/{topicId}", h.showTopicView).Methods(http.MethodGet)
	router.HandleFunc("/topics/{topicId}/edit", h.protect(h.showEditTopicView)).Methods(http.MethodGet)
	router.HandleFunc("/topics/{topicId}/update", h.protect(h.updateTopic)).Methods(http.MethodPost)

	router.HandleFunc("/posts/save", h.protect(h.savePost)).Methods(http.MethodPost)
	//router.HandleFunc("/posts/{postId}", h.showPostView).Methods(http.MethodGet)
	router.HandleFunc("/posts/{postId}/edit", h.protect(h.showEditPostView)).Methods(http.MethodGet)
	router.HandleFunc("/posts/{postId}/update", h.protect(h.updatePost)).Methods(http.MethodPost)
	router.HandleFunc("/posts/{postId}/remove", h.protect(h.handleRemovePost))
	//router.HandleFunc("/posts/{postId}/reply", h.protect(h.savePostReply)).Methods(http.MethodPost)

	// Admin
	router.HandleFunc("/admin", h.admin(h.showAdminView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/boards", h.admin(h.showAdminBoardsView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/boards/new", h.admin(h.showNewBoardView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/boards/save", h.admin(h.saveBoard)).Methods(http.MethodPost)
	router.HandleFunc("/admin/boards/{boardId}/edit", h.admin(h.showEditBoardView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/boards/{boardId}/update", h.admin(h.updateBoard)).Methods(http.MethodPost)
	router.HandleFunc("/admin/users", h.admin(h.showAdminUsersView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/users/{name}/edit", h.admin(h.showEditUserView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/users/{name}/update", h.admin(h.updateUserAdmin)).Methods(http.MethodPost)
	router.HandleFunc("/admin/settings/edit", h.admin(h.showAdminSettingsView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/settings/update", h.admin(h.updateSettingsAdmin)).Methods(http.MethodPost)
	router.HandleFunc("/admin/keys", h.admin(h.showKeysView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/keys/save", h.admin(h.saveKey)).Methods(http.MethodPost)
	router.HandleFunc("/reset-password", h.showResetPasswordView).Methods(http.MethodGet)
	router.HandleFunc("/reset-password", h.updatePassword).Methods(http.MethodPost)

	router.HandleFunc("/admin/forums", h.admin(h.showAdminForumsView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/forums/new", h.admin(h.showNewForumView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/forums/save", h.admin(h.saveForum)).Methods(http.MethodPost)
	router.HandleFunc("/admin/forums/{forumId}/edit", h.admin(h.showEditForumView)).Methods(http.MethodGet)
	router.HandleFunc("/admin/forums/{forumId}/update", h.admin(h.updateForum)).Methods(http.MethodPost)

	router.HandleFunc("/style.css", h.showStylesheet).Methods(http.MethodGet)

	// Pagination
	//router.HandleFunc("/page/{nb}", h.showPageNumber).Methods(http.MethodGet)

	// Replies
	//router.HandleFunc("/replies/{replyId}", h.protect(h.showReplyView)).Methods(http.MethodGet)
	//router.HandleFunc("/replies/{replyId}/save", h.protect(h.saveReplyReply)).Methods(http.MethodPost)
	//router.HandleFunc("/replies/{replyId}/edit", h.protect(h.showEditReplyView)).Methods(http.MethodGet)
	//router.HandleFunc("/replies/{replyId}/update", h.protect(h.updateReply)).Methods(http.MethodPost)
	//router.HandleFunc("/replies/{replyId}/remove", h.protect(h.handleRemoveReply))

	// Notifications
	//router.HandleFunc("/notifications", h.protect(h.showNotificationsView)).Methods(http.MethodGet)
	//router.HandleFunc("/notifications/{notificationId}/mark-read", h.protect(h.markRead)).Methods(http.MethodGet)
	//router.HandleFunc("/notifications/mark-all-read", h.protect(h.markAllRead)).Methods(http.MethodGet)

	// User
	router.HandleFunc("/~{userId}", h.showUserPostsView).Methods(http.MethodGet)
	router.HandleFunc("/account", h.protect(h.showAccountView)).Methods(http.MethodGet)
	router.HandleFunc("/save-account", h.protect(h.saveAccount)).Methods(http.MethodPost)

	// Index
	router.HandleFunc("/", h.showIndexView).Name("index").Methods(http.MethodGet)

	return router, nil
}
