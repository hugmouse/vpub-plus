package handler

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminView(w http.ResponseWriter, r *http.Request, user model.User) {
	h.renderLayout(w, "admin", nil, user)
}

func (h *Handler) showEditUserView(w http.ResponseWriter, r *http.Request, user model.User) {
	u, err := h.storage.UserByName(mux.Vars(r)["name"])
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "admin_user_edit", map[string]interface{}{
		"user": u,
		"form": form.AdminUserForm{
			Username: u.Name,
			About:    u.About,
		},
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showAdminUsersView(w http.ResponseWriter, r *http.Request, user model.User) {
	users, err := h.storage.Users()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "admin_user", map[string]interface{}{
		"users": users,
	}, user)
}

func (h *Handler) showAdminBoardsView(w http.ResponseWriter, r *http.Request, user model.User) {
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	hasForums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	forums := forumFromBoards(boards)
	h.renderLayout(w, "admin_board", map[string]interface{}{
		"hasForums": hasForums,
		"forums":    forums,
	}, user)
}

func (h *Handler) showAdminForumsView(w http.ResponseWriter, r *http.Request, user model.User) {
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "admin_forum", map[string]interface{}{
		"forums": forums,
	}, user)
}

func (h *Handler) showKeysView(w http.ResponseWriter, r *http.Request, user model.User) {
	keys, err := h.storage.Keys()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "admin_keys", map[string]interface{}{
		"keys":           keys,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showAdminSettingsView(w http.ResponseWriter, r *http.Request, user model.User) {
	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}
	settingsForm := form.SettingsForm{
		Name: settings.Name,
		Css:  settings.Css,
	}
	h.renderLayout(w, "admin_settings_edit", map[string]interface{}{
		"form":           settingsForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) updateSettingsAdmin(w http.ResponseWriter, r *http.Request, user model.User) {
	settingsForm := form.NewSettingsForm(r)
	var settings model.Settings
	settings.Name = settingsForm.Name
	settings.Css = settingsForm.Css
	settings.Footer = settingsForm.Footer
	if err := h.storage.UpdateSettings(settings); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handler) showNewBoardView(w http.ResponseWriter, r *http.Request, user model.User) {
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.BoardForm{
		Forums: forums,
	}
	h.renderLayout(w, "admin_board_create", map[string]interface{}{
		"form":           boardForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showNewForumView(w http.ResponseWriter, r *http.Request, user model.User) {
	forumForm := form.ForumForm{}
	h.renderLayout(w, "admin_forum_create", map[string]interface{}{
		"form":           forumForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showEditBoardView(w http.ResponseWriter, r *http.Request, user model.User) {
	board, err := h.storage.BoardById(RouteInt64Param(r, "boardId"))
	if err != nil {
		serverError(w, err)
		return
	}
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.BoardForm{
		Name:        board.Name,
		Description: board.Description,
		Position:    board.Position,
		ForumId:     board.Forum.Id,
		Forums:      forums,
		IsLocked:    board.IsLocked,
	}
	h.renderLayout(w, "admin_board_edit", map[string]interface{}{
		"form":           boardForm,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showEditForumView(w http.ResponseWriter, r *http.Request, user model.User) {
	forum, err := h.storage.ForumById(RouteInt64Param(r, "forumId"))
	if err != nil {
		serverError(w, err)
		return
	}
	forumForm := form.ForumForm{
		Name:     forum.Name,
		Position: forum.Position,
	}
	h.renderLayout(w, "admin_forum_edit", map[string]interface{}{
		"forum":          forum,
		"form":           forumForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) updateForum(w http.ResponseWriter, r *http.Request, user model.User) {
	id := RouteInt64Param(r, "forumId")
	forum, err := h.storage.ForumById(id)
	if err != nil {
		serverError(w, err)
		return
	}
	forumForm := form.NewForumForm(r)
	forum.Name = forumForm.Name
	forum.Position = forumForm.Position
	if err := h.storage.UpdateForum(forum); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}

func (h *Handler) updateBoard(w http.ResponseWriter, r *http.Request, user model.User) {
	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.NewBoardForm(r)
	board.Name = boardForm.Name
	board.Description = boardForm.Description
	board.Position = boardForm.Position
	board.Forum.Id = boardForm.ForumId
	board.IsLocked = boardForm.IsLocked
	if err := h.storage.UpdateBoard(board); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}

func (h *Handler) updateUserAdmin(w http.ResponseWriter, r *http.Request, user model.User) {
	userForm := form.NewAdminUserForm(r)
	user.Name = userForm.Username
	user.About = userForm.About
	if err := h.storage.UpdateUser(user); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func (h *Handler) saveBoard(w http.ResponseWriter, r *http.Request, user model.User) {
	boardForm := form.NewBoardForm(r)
	board := model.Board{
		Name:        boardForm.Name,
		Description: boardForm.Description,
		Position:    boardForm.Position,
		IsLocked:    boardForm.IsLocked,
		Forum:       model.Forum{Id: boardForm.ForumId},
	}
	_, err := h.storage.CreateBoard(board)
	if err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}

func (h *Handler) saveForum(w http.ResponseWriter, r *http.Request, user model.User) {
	forumForm := form.NewForumForm(r)
	forum := model.Forum{
		Name:     forumForm.Name,
		Position: forumForm.Position,
	}
	_, err := h.storage.CreateForum(forum)
	if err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}

func (h *Handler) saveKey(w http.ResponseWriter, r *http.Request, user model.User) {
	if err := h.storage.CreateKey(); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/keys", http.StatusFound)
}
