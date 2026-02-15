package handler

import (
	"context"
	_ "expvar"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
	"vpub/model"
	"vpub/storage"
	"vpub/syntax"
	"vpub/syntax/renderers/customBlackfriday"
	"vpub/syntax/renderers/vanilla"
	"vpub/web/handler/request"
	"vpub/web/session"

	"github.com/gorilla/mux"
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

var errorTemplate = template.Must(template.New("error").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<link rel="stylesheet" href="/style.css">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Server Error</title>
</head>
<body>
	<div class="search-result">
		<h1>Server Error</h1>
		<pre><code>{{ .ErrorMessage }}</code></pre>
	</div>
</body>
</html>
`))

func serverError(w http.ResponseWriter, err error) {
	log.Println("[server error]", err)
	w.WriteHeader(http.StatusInternalServerError)
	if tplErr := errorTemplate.Execute(w, struct {
		ErrorMessage string
	}{
		ErrorMessage: err.Error(),
	}); tplErr != nil {
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
	}
}

func (h *Handler) getCachedSettings() (model.Settings, error) {
	h.settingsCacheMutex.RLock()
	if time.Since(h.settingsCacheTime) < h.settingsCacheTTL && h.settingsCache != nil {
		cached := *h.settingsCache
		h.settingsCacheMutex.RUnlock()
		return cached, nil
	}
	h.settingsCacheMutex.RUnlock()

	// Cache miss, fetch from DB
	h.settingsCacheMutex.Lock()
	defer h.settingsCacheMutex.Unlock()

	settings, err := h.storage.Settings()
	if err != nil {
		return model.Settings{}, err
	}

	if settings.SettingsCacheTTL > 0 {
		h.settingsCacheTTL = time.Duration(settings.SettingsCacheTTL) * time.Second
	} else {
		h.settingsCacheTTL = 0
	}

	h.settingsCache = &settings
	h.settingsCacheTime = time.Now()
	return settings, nil
}

func (h *Handler) invalidateSettingsCache() {
	h.settingsCacheMutex.Lock()
	defer h.settingsCacheMutex.Unlock()
	h.settingsCache = nil
}

func (h *Handler) handleSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		settings, err := h.getCachedSettings()
		if err != nil {
			serverError(w, err)
			return
		}
		ctx = context.WithValue(ctx, request.SettingsKey, settings)
		user, _session, err := h.session.GetUser(r)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, request.SessionKey, _session)
		ctx = context.WithValue(ctx, request.UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type Handler struct {
	session             *session.Manager
	mux                 *mux.Router
	storage             *storage.Storage
	currentRenderEngine *syntax.Renderer
	renderRegistry      *syntax.RenderEngineRegistry
	imageProxy          *ImageProxyHandler
	settingsCache       *model.Settings
	settingsCacheMutex  sync.RWMutex
	settingsCacheTime   time.Time
	settingsCacheTTL    time.Duration
}

type ImageProxyHandler struct {
	httpClient   *http.Client
	cachedImages map[string]CachedImage
	cacheMutex   sync.RWMutex
}

type CachedImage struct {
	lastUpdate  time.Time
	value       interface{}
	contentType string
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
			forum.ID = board.Forum.ID
		} else if board.Forum.ID != forum.ID {
			forums = append(forums, forum)
			forum = model.Forum{Name: board.Forum.Name, ID: board.Forum.ID}
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

	renderRegistry := syntax.NewRegistry()

	err := renderRegistry.Register("blackfriday", &customBlackfriday.BlackfridayRenderer{})
	if err != nil {
		return nil, err
	}
	err = renderRegistry.Register("vanilla", &vanilla.VanillaRenderer{})
	if err != nil {
		return nil, err
	}

	defaultRenderEngine, err := renderRegistry.Get("blackfriday")
	if err != nil {
		return nil, err
	}

	h := &Handler{
		session:             s,
		mux:                 router,
		storage:             data,
		currentRenderEngine: &defaultRenderEngine,
		renderRegistry:      renderRegistry,
		settingsCacheTTL:    30 * time.Second,
	}

	router.Use(h.handleSessionMiddleware)
	h.initTpl()

	handlerForImageProxy := &ImageProxyHandler{
		cachedImages: make(map[string]CachedImage),
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 3 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 3 * time.Second,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}

	h.imageProxy = handlerForImageProxy

	// Adds pprof to /debug/pprof route,
	// see "debug_handlers.go" for more info
	registerDebugHandlers(router)

	// Static assets
	router.HandleFunc("/style.css", h.showStylesheet).Methods(http.MethodGet)
	router.HandleFunc("/js/{filename}", h.showJS).Methods(http.MethodGet)
	//router.HandleFunc("/favicon.ico", h.showFavicon).Name("favicon").Methods(http.MethodGet)

	// Forum views
	publicSubRouter := router.PathPrefix("/").Subrouter()

	// Proxy
	publicSubRouter.HandleFunc("/image-proxy", handlerForImageProxy.imageProxyHandler).Methods(http.MethodGet)

	// Search
	publicSubRouter.HandleFunc("/search", h.searchShow).Methods(http.MethodGet)

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

	adminSubRouter.HandleFunc("/image-proxy", h.showAdminImageCache).Methods(http.MethodGet)
	adminSubRouter.HandleFunc("/image-proxy/remove", h.removeAdminImageCache).Methods(http.MethodPost)

	return router, nil
}
