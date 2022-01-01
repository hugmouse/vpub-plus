package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"vpub/config"
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

type ProtectedFunc func(http.ResponseWriter, *http.Request, string)

type Handler struct {
	session *session.Session
	host    string
	env     string
	css     []byte
	mux     *mux.Router
	storage *storage.Storage
	title   string
	motd    []byte
	topics  []string
	perPage int
}

func (h *Handler) protect(fn ProtectedFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.session.Get(r)
		if err != nil || user == "" {
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
		topics:  cfg.Topics,
		perPage: cfg.PerPage,
	}
	h.initTpl()

	// Read and cache css
	cssFile, _ := os.Open(cfg.CSSFile)
	b, err := io.ReadAll(cssFile)
	if err != nil {
		fmt.Println("Couldn't read CSS file. Set CSS_FILE. Value: ", cfg.CSSFile)
	}
	h.css = b
	// Read and cache motd
	motdFile, _ := os.Open(cfg.MOTDFile)
	b, err = io.ReadAll(motdFile)
	if err != nil {
		fmt.Println("Couldn't read MOTD file. Set MOTD_FILE. Value:", cfg.MOTDFile)
	}
	h.motd = b
	h.title = cfg.Title

	// Static assets
	router.HandleFunc("/style.css", h.showStylesheet).Methods(http.MethodGet)
	//router.HandleFunc("/manual", h.showManual).Name("manual").Methods(http.MethodGet)
	//router.HandleFunc("/favicon.ico", h.showFavicon).Name("favicon").Methods(http.MethodGet)
	router.HandleFunc("/feed.atom", h.showFeedView).Methods(http.MethodGet)

	// Auth
	router.HandleFunc("/login", h.showLoginView).Methods(http.MethodGet)
	router.HandleFunc("/login", h.checkLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.showRegisterView).Methods(http.MethodGet)
	router.HandleFunc("/register", h.register).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.logout).Methods(http.MethodGet)

	// Topics
	router.HandleFunc("/topics/{topic}", h.showTopicView).Methods(http.MethodGet)
	router.HandleFunc("/topics/{topic}/feed.atom", h.showFeedViewTopic).Methods(http.MethodGet)

	// Posts
	//router.HandleFunc("/posts", h.showPostsView).Name("posts").Methods(http.MethodGet)
	router.HandleFunc("/posts/new", h.protect(h.showNewPostView)).Methods(http.MethodGet)
	router.HandleFunc("/posts/save", h.protect(h.savePost)).Methods(http.MethodPost)
	router.HandleFunc("/posts/{postId}", h.showPostView).Methods(http.MethodGet)
	router.HandleFunc("/posts/{postId}/edit", h.protect(h.showEditPostView)).Methods(http.MethodGet)
	router.HandleFunc("/posts/{postId}/update", h.protect(h.updatePost)).Methods(http.MethodPost)
	router.HandleFunc("/posts/{postId}/remove", h.protect(h.handleRemovePost))
	router.HandleFunc("/posts/{postId}/reply", h.protect(h.savePostReply)).Name("savePostReply").Methods(http.MethodPost)

	// Pagination
	router.HandleFunc("/page/{nb}", h.showPageNumber).Methods(http.MethodGet)

	// Replies
	router.HandleFunc("/replies/{replyId}", h.protect(h.showReplyView)).Methods(http.MethodGet)
	router.HandleFunc("/replies/{replyId}/save", h.protect(h.saveReplyReply)).Methods(http.MethodPost)
	router.HandleFunc("/replies/{replyId}/edit", h.protect(h.showEditReplyView)).Methods(http.MethodGet)
	router.HandleFunc("/replies/{replyId}/update", h.protect(h.updateReply)).Methods(http.MethodPost)
	router.HandleFunc("/replies/{replyId}/remove", h.protect(h.handleRemoveReply)).Name("removeReply")
	//
	//// Notifications
	router.HandleFunc("/notifications", h.protect(h.showNotificationsView)).Methods(http.MethodGet)
	router.HandleFunc("/notifications/{notificationId}/mark-read", h.protect(h.markRead)).Methods(http.MethodPost)
	router.HandleFunc("/notifications/mark-all-read", h.protect(h.markAllRead)).Methods(http.MethodPost)
	//
	//// User
	router.HandleFunc("/~{userId}", h.showUserPostsView).Methods(http.MethodGet)
	//router.HandleFunc("/account", h.protect(h.showAccountView)).Name("account").Methods(http.MethodGet)
	//router.HandleFunc("/patrons", h.showUserListView).Name("patrons").Methods(http.MethodGet)
	//router.HandleFunc("/save-about", h.protect(h.saveAbout)).Name("saveAbout").Methods(http.MethodPost)
	//router.HandleFunc("/site", h.protect(h.showSiteView)).Name("site").Methods(http.MethodGet)
	//router.HandleFunc("/theme", h.protect(h.showEditThemeView)).Name("editTheme").Methods(http.MethodGet)
	//router.HandleFunc("/theme/update", h.protect(h.updateTheme)).Name("updateTheme").Methods(http.MethodPost)
	//
	//// Index
	router.HandleFunc("/", h.showIndexView).Name("index").Methods(http.MethodGet)
	//
	//engine, err := template.New(env, host, h)
	//if err != nil {
	//	return router, err
	//}
	//
	//h.tpl = engine
	//
	//mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	r.Body = http.MaxBytesReader(w, r.Body, 3<<20)
	//	router.ServeHTTP(w, r)
	//})

	return router, nil
	//return mux, err
}
