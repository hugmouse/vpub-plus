package handler

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminView(w http.ResponseWriter, r *http.Request) {
	h.renderLayout(w, r, "admin", nil)
}

func (h *Handler) showEditUserView(w http.ResponseWriter, r *http.Request) {
	u, err := h.storage.UserByName(mux.Vars(r)["name"])
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, r, "admin_user_edit", map[string]interface{}{
		"user": u,
		"form": form.AdminUserForm{
			Username: u.Name,
			About:    u.About,
		},
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) showAdminUsersView(w http.ResponseWriter, r *http.Request) {
	users, err := h.storage.Users()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, r, "admin_user", map[string]interface{}{
		"users": users,
	})
}

func (h *Handler) showAdminBoardsView(w http.ResponseWriter, r *http.Request) {
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
	h.renderLayout(w, r, "admin_board", map[string]interface{}{
		"hasForums": hasForums,
		"forums":    forums,
	})
}

func (h *Handler) showAdminForumsView(w http.ResponseWriter, r *http.Request) {
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, r, "admin_forum", map[string]interface{}{
		"forums": forums,
	})
}

func (h *Handler) showKeysView(w http.ResponseWriter, r *http.Request) {
	keys, err := h.storage.Keys()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, r, "admin_keys", map[string]interface{}{
		"keys":           keys,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) showAdminSettingsView(w http.ResponseWriter, r *http.Request) {
	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}
	settingsForm := form.SettingsForm{
		Name:    settings.Name,
		Css:     settings.Css,
		Footer:  settings.Footer,
		PerPage: settings.PerPage,
	}
	h.renderLayout(w, r, "admin_settings_edit", map[string]interface{}{
		"form":           settingsForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) updateSettingsAdmin(w http.ResponseWriter, r *http.Request) {
	settingsForm := form.NewSettingsForm(r)
	var settings model.Settings
	settings.Name = settingsForm.Name
	settings.Css = settingsForm.Css
	settings.Footer = settingsForm.Footer
	settings.PerPage = settingsForm.PerPage
	if err := h.storage.UpdateSettings(settings); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handler) showNewBoardView(w http.ResponseWriter, r *http.Request) {
	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.BoardForm{
		Forums: forums,
	}
	h.renderLayout(w, r, "admin_board_create", map[string]interface{}{
		"form":           boardForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) showNewForumView(w http.ResponseWriter, r *http.Request) {
	forumForm := form.ForumForm{}
	h.renderLayout(w, r, "admin_forum_create", map[string]interface{}{
		"form":           forumForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) showEditBoardView(w http.ResponseWriter, r *http.Request) {
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
	h.renderLayout(w, r, "admin_board_edit", map[string]interface{}{
		"form":           boardForm,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) showEditForumView(w http.ResponseWriter, r *http.Request) {
	forum, err := h.storage.ForumById(RouteInt64Param(r, "forumId"))
	if err != nil {
		serverError(w, err)
		return
	}
	forumForm := form.ForumForm{
		Name:     forum.Name,
		Position: forum.Position,
		IsLocked: forum.IsLocked,
	}
	h.renderLayout(w, r, "admin_forum_edit", map[string]interface{}{
		"forum":          forum,
		"form":           forumForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func (h *Handler) updateForum(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "forumId")
	forum, err := h.storage.ForumById(id)
	if err != nil {
		serverError(w, err)
		return
	}
	forumForm := form.NewForumForm(r)
	forum.Name = forumForm.Name
	forum.Position = forumForm.Position
	forum.IsLocked = forumForm.IsLocked
	if err := h.storage.UpdateForum(forum); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}

func (h *Handler) updateBoard(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) updateUserAdmin(w http.ResponseWriter, r *http.Request) {
	user, _ := h.session.Get(r)
	userForm := form.NewAdminUserForm(r)
	user.Name = userForm.Username
	user.About = userForm.About
	if err := h.storage.UpdateUser(user); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func (h *Handler) saveBoard(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) saveForum(w http.ResponseWriter, r *http.Request) {
	forumForm := form.NewForumForm(r)
	forum := model.Forum{
		Name:     forumForm.Name,
		Position: forumForm.Position,
		IsLocked: forumForm.IsLocked,
	}
	_, err := h.storage.CreateForum(forum)
	if err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}

func (h *Handler) saveKey(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.CreateKey(); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/keys", http.StatusFound)
}
