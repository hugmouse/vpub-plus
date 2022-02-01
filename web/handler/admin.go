package handler

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
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
	v := NewView(w, r, "admin_board")
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
	v.Set("hasForums", hasForums)
	v.Set("forums", forums)
	v.Render()
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
	if err := h.storage.UpdateSettings(*settingsForm.Merge(&model.Settings{})); err != nil {
		serverError(w, err)
		return
	}
	session := request.GetSessionContextKey(r)
	session.FlashInfo("Settings updated")
	session.Save(r, w)
	http.Redirect(w, r, "/admin/settings/edit", http.StatusFound)
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
	v := NewView(w, r, "admin_board_create")
	v.Set("form", boardForm)
	v.Render()
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
	v := NewView(w, r, "admin_board_edit")
	v.Set("form", boardForm)
	v.Set("board", board)
	v.Render()
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
		notFound(w)
		return
	}

	// 1. Build the form object
	forumForm := form.NewForumForm(r)

	// 2. Prepare the view
	v := NewView(w, r, "admin_forum_edit")
	v.Set("forum", forum)
	v.Set("form", forumForm)

	// 3. Create request
	forumRequest := model.ForumRequest{
		Name:     forumForm.Name,
		Position: forumForm.Position,
		IsLocked: forumForm.IsLocked,
	}

	// 4. Validate request
	if err := validator.ValidateForumModification(h.storage, id, forumRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	// 5. Process request
	if err := h.storage.UpdateForum(forumRequest.Patch(forum)); err != nil {
		v.Set("errorMessage", "Unable to update forum")
		serverError(w, err)
		return
	}

	// 6. Happy path
	session := request.GetSessionContextKey(r)
	session.FlashInfo("Forum updated")
	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/admin/forums/%d/edit", id), http.StatusFound)
}

func (h *Handler) updateBoard(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "boardId")

	board, err := h.storage.BoardById(id)
	if err != nil {
		serverError(w, err)
		return
	}

	boardForm := form.NewBoardForm(r)

	v := NewView(w, r, "admin_board_edit")
	v.Set("board", board)
	v.Set("form", boardForm)

	if err := boardForm.Validate(); err != nil {
		v.Set("errorMessage", err.Error())
		forums, err := h.storage.Forums()
		if err != nil {
			serverError(w, err)
			return
		}
		boardForm.Forums = forums
		v.Render()
		return
	}

	if err := h.storage.UpdateBoard(*boardForm.Merge(&board)); err != nil {
		serverError(w, err)
		return
	}

	session := request.GetSessionContextKey(r)
	session.FlashInfo("Board updated")
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/admin/boards/%d/edit", id), http.StatusFound)
}

func (h *Handler) updateUserAdmin(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)
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

	v := NewView(w, r, "admin_board_create")
	v.Set("form", boardForm)

	if err := boardForm.Validate(); err != nil {
		v.Set("errorMessage", err.Error())
		forums, err := h.storage.Forums()
		if err != nil {
			serverError(w, err)
			return
		}
		boardForm.Forums = forums
		v.Render()
		return
	}

	_, err := h.storage.CreateBoard(*boardForm.Merge(&model.Board{}))
	if err != nil {
		serverError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}

func (h *Handler) saveForum(w http.ResponseWriter, r *http.Request) {
	// 1. Build the form object
	forumForm := form.NewForumForm(r)

	// 2. Prepare the view
	v := NewView(w, r, "admin_forum_create")
	v.Set("form", forumForm)

	// 3. Create request
	forumRequest := model.ForumRequest{
		Name:     forumForm.Name,
		Position: forumForm.Position,
		IsLocked: forumForm.IsLocked,
	}

	// 4. Validate request
	if err := validator.ValidateForumCreation(h.storage, forumRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	// 5. Process the request
	if _, err := h.storage.CreateForum(forumRequest); err != nil {
		v.Set("errorMessage", "Unable to create forum")
		serverError(w, err)
		return
	}

	// 6. Happy path
	http.Redirect(w, r, "/admin/forums", http.StatusFound)
}

func (h *Handler) saveKey(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.CreateKey(); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/keys", http.StatusFound)
}
