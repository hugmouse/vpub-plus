package handler

import (
	"context"
	_ "expvar"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"vpub/model"
	"vpub/storage"
	"vpub/syntax"
	"vpub/syntax/renderers/customBlackfriday"
	"vpub/syntax/renderers/vanilla"
	"vpub/web/handler/request"
	"vpub/web/session"
)

func RouteInt64Param(r *http.Request, param string) int64 {
	value, err := strconv.ParseInt(r.PathValue(param), 10, 64)
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

func (h *Handler) handleSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		settings, err := h.storage.Settings()
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
	storage             *storage.Storage
	currentRenderEngine *syntax.Renderer
	renderRegistry      *syntax.RenderEngineRegistry
	imageProxy          *ImageProxyHandler
	setupComplete       atomic.Bool
}

type ImageProxyHandler struct {
	httpClient   *http.Client
	cachedImages *lruCache
}

type lruCache struct {
	mu       sync.Mutex
	items    map[string]*lruItem
	capacity int
}

type lruItem struct {
	value       []byte
	contentType string
	lastUpdate  time.Time
	lastAccess  time.Time
}

func newLRUCache(capacity int) *lruCache {
	if capacity <= 0 {
		capacity = 100
	}
	return &lruCache{
		items:    make(map[string]*lruItem),
		capacity: capacity,
	}
}

func (c *lruCache) Get(key string) (lruItem, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return lruItem{}, false
	}
	item.lastAccess = time.Now()
	return *item, true
}

func (c *lruCache) Set(key string, item *lruItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.items[key]; !exists && len(c.items) >= c.capacity {
		c.evictOldest()
	}
	item.lastAccess = time.Now()
	c.items[key] = item
}

func (c *lruCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*lruItem)
}

func (c *lruCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true
	for k, v := range c.items {
		if first || v.lastAccess.Before(oldestTime) {
			oldestKey = k
			oldestTime = v.lastAccess
			first = false
		}
	}
	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

func (c *lruCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

func (c *lruCache) Items() map[string]lruItem {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make(map[string]lruItem)
	for k, v := range c.items {
		result[k] = *v
	}
	return result
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

// setupGuard redirects every request to /setup until the initial onboarding has
// been completed. Static assets, the health check and the setup endpoints
// themselves are always allowed through.
func (h *Handler) setupGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.setupComplete.Load() {
			next.ServeHTTP(w, r)
			return
		}
		adminExists, err := h.storage.HasAdmin()
		if err != nil {
			serverError(w, err)
			return
		}
		if adminExists {
			h.setupComplete.Store(true)
			next.ServeHTTP(w, r)
			return
		}
		path := r.URL.Path
		if path == "/setup" || strings.HasPrefix(path, "/style.css") ||
			strings.HasPrefix(path, "/js/") || path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/setup", http.StatusFound)
	})
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

func New(data *storage.Storage, s *session.Manager, csrfSecure bool) (http.Handler, error) {
	mux := http.NewServeMux()

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
		storage:             data,
		currentRenderEngine: &defaultRenderEngine,
		renderRegistry:      renderRegistry,
	}

	h.initTpl()

	handlerForImageProxy := &ImageProxyHandler{
		cachedImages: newLRUCache(200),
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
				return validateImageProxyURL(req.URL)
			},
		},
	}

	h.imageProxy = handlerForImageProxy

	registerDebugHandlers(mux)

	// Setup (onboarding; when user just installed vpub-plus)
	mux.HandleFunc("GET /setup", h.showSetupView)
	mux.HandleFunc("POST /setup", h.setup)

	// Static assets
	mux.HandleFunc("GET /style.css", h.showStylesheet)
	mux.HandleFunc("GET /js/{filename}", h.showJS)

	// Health check
	mux.HandleFunc("GET /health", h.healthCheck)

	// Proxy
	mux.HandleFunc("GET /image-proxy", handlerForImageProxy.imageProxyHandler)

	// Search
	mux.HandleFunc("GET /search", h.searchShow)

	// Auth
	mux.HandleFunc("GET /login", h.showLoginView)
	mux.HandleFunc("POST /login", h.checkLogin)
	mux.HandleFunc("GET /register", h.showRegisterView)
	mux.HandleFunc("POST /register", h.register)
	mux.HandleFunc("GET /logout", h.logout)

	// Feed
	mux.HandleFunc("GET /feed.atom", h.showFeed)

	// Forums
	mux.HandleFunc("GET /forums/{forumId}", h.showForumView)

	// Boards
	mux.HandleFunc("GET /boards/{boardId}", h.showBoardView)
	mux.HandleFunc("GET /boards/{boardId}/feed.atom", h.showBoardFeed)
	mux.HandleFunc("GET /boards/{boardId}/new-topic", h.protect(h.showCreateTopicView))
	mux.HandleFunc("POST /boards/{boardId}/save-topic", h.protect(h.saveTopic))
	mux.HandleFunc("GET /boards/{boardId}/newest", h.showNewestBoardView)

	// Topic
	mux.HandleFunc("GET /topics/{topicId}", h.showTopicView)
	mux.HandleFunc("GET /topics/{topicId}/feed.atom", h.showTopicFeed)
	mux.HandleFunc("GET /topics/{topicId}/edit", h.protect(h.showEditTopicView))
	mux.HandleFunc("POST /topics/{topicId}/update", h.updateTopic)
	mux.HandleFunc("GET /topics/{topicId}/newest", h.showNewestTopicView)

	// Post
	mux.HandleFunc("POST /posts/save", h.protect(h.savePost))
	mux.HandleFunc("GET /posts/{postId}/edit", h.protect(h.showEditPostView))
	mux.HandleFunc("POST /posts/{postId}/update", h.protect(h.updatePost))
	mux.HandleFunc("GET /posts/{postId}/remove", h.protect(h.removePost))
	mux.HandleFunc("POST /posts/{postId}/remove", h.protect(h.removePost))
	mux.HandleFunc("GET /posts", h.showPostListView)

	// Account
	mux.HandleFunc("GET /account", h.protect(h.showAccountEditPage))
	mux.HandleFunc("POST /update-account", h.protect(h.updateAccount))
	mux.HandleFunc("GET /reset-password", h.showResetPasswordView)
	mux.HandleFunc("POST /reset-password", h.updatePassword)

	// Users
	mux.HandleFunc("GET /users/{userId}", h.showUserView)

	// Index — {$} means exact match on "/" only
	mux.HandleFunc("GET /{$}", h.showIndexView)

	// Admin routes — each individually wrapped with h.admin()
	mux.HandleFunc("GET /admin/", h.admin(h.showAdminView))
	mux.HandleFunc("GET /admin/boards", h.admin(h.showAdminBoardsView))
	mux.HandleFunc("GET /admin/boards/new", h.admin(h.showAdminCreateBoardView))
	mux.HandleFunc("POST /admin/boards/save", h.admin(h.saveAdminBoard))
	mux.HandleFunc("GET /admin/boards/{boardId}/edit", h.admin(h.showAdminEditBoardView))
	mux.HandleFunc("POST /admin/boards/{boardId}/update", h.admin(h.updateAdminBoard))
	mux.HandleFunc("GET /admin/boards/{boardId}/remove", h.admin(h.showAdminRemoveBoardView))
	mux.HandleFunc("POST /admin/boards/{boardId}/remove", h.admin(h.removeAdminBoard))

	mux.HandleFunc("GET /admin/users", h.admin(h.showAdminUserListView))
	mux.HandleFunc("GET /admin/users/{userId}/edit", h.admin(h.showAdminEditUserView))
	mux.HandleFunc("POST /admin/users/{userId}/update", h.admin(h.updateAdminUser))
	mux.HandleFunc("GET /admin/users/{userId}/remove", h.admin(h.showAdminRemoveUserView))
	mux.HandleFunc("POST /admin/users/{userId}/remove", h.admin(h.removeAdminUser))

	mux.HandleFunc("GET /admin/settings/edit", h.admin(h.showAdminSettingsView))
	mux.HandleFunc("POST /admin/settings/update", h.admin(h.updateAdminSettings))

	mux.HandleFunc("GET /admin/keys", h.admin(h.showAdminKeyListView))
	mux.HandleFunc("POST /admin/keys/save", h.admin(h.saveAdminKey))
	mux.HandleFunc("POST /admin/keys/{keyId}/remove", h.admin(h.removeAdminKey))

	mux.HandleFunc("GET /admin/forums", h.admin(h.showAdminForumsView))
	mux.HandleFunc("GET /admin/forums/new", h.admin(h.showAdminCreateForumView))
	mux.HandleFunc("POST /admin/forums/save", h.admin(h.saveAdminForum))
	mux.HandleFunc("GET /admin/forums/{forumId}/edit", h.admin(h.showAdminEditForumView))
	mux.HandleFunc("POST /admin/forums/{forumId}/update", h.admin(h.updateAdminForum))
	mux.HandleFunc("GET /admin/forums/{forumId}/remove", h.admin(h.showAdminRemoveForumView))
	mux.HandleFunc("POST /admin/forums/{forumId}/remove", h.admin(h.removeAdminForum))

	mux.HandleFunc("GET /admin/image-proxy", h.admin(h.showAdminImageCache))
	mux.HandleFunc("POST /admin/image-proxy/remove", h.admin(h.removeAdminImageCache))

	var handler http.Handler = mux
	handler = h.handleSessionMiddleware(handler)
	handler = h.setupGuard(handler)
	handler = newCSRFMiddleware(csrfSecure)(handler)
	return handler, nil
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Ping(); err != nil {
		http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("health check: failed to write response: %v", err)
	}
}
