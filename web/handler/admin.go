package handler

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
)

func (h *Handler) showAdminView(w http.ResponseWriter, r *http.Request, user string) {
	h.renderLayout(w, "admin", nil, user)
}

func (h *Handler) showEditUserView(w http.ResponseWriter, r *http.Request, user string) {
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

func (h *Handler) showAdminUsersView(w http.ResponseWriter, r *http.Request, user string) {
	users, err := h.storage.Users()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "admin_user", map[string]interface{}{
		"users": users,
	}, user)
}

func (h *Handler) showAdminBoardsView(w http.ResponseWriter, r *http.Request, user string) {
	boards, err := h.storage.Boards()
	if err != nil {
		serverError(w, err)
		return
	}
	h.renderLayout(w, "admin_board", map[string]interface{}{
		"boards": boards,
	}, user)
}

func (h *Handler) showNewBoardView(w http.ResponseWriter, r *http.Request, user string) {
	boardForm := form.BoardForm{}
	h.renderLayout(w, "admin_board_create", map[string]interface{}{
		"form":           boardForm,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) showEditBoardView(w http.ResponseWriter, r *http.Request, user string) {
	board, err := h.storage.BoardById(RouteInt64Param(r, "boardId"))
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.BoardForm{
		Name:        board.Name,
		Description: board.Description,
	}
	h.renderLayout(w, "admin_board_edit", map[string]interface{}{
		"form":           boardForm,
		"board":          board,
		csrf.TemplateTag: csrf.TemplateField(r),
	}, user)
}

func (h *Handler) updateBoard(w http.ResponseWriter, r *http.Request, user string) {
	id := RouteInt64Param(r, "boardId")
	board, err := h.storage.BoardById(id)
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm := form.NewBoardForm(r)
	board.Name = boardForm.Name
	board.Description = boardForm.Description
	if err := h.storage.UpdateBoard(board); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}

func (h *Handler) updateUserAdmin(w http.ResponseWriter, r *http.Request, name string) {
	userForm := form.NewAdminUserForm(r)
	user := model.User{
		Name:  userForm.Username,
		About: userForm.About,
	}
	if err := h.storage.UpdateUser(name, user); err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func (h *Handler) saveBoard(w http.ResponseWriter, r *http.Request, user string) {
	boardForm := form.NewBoardForm(r)
	board := model.Board{
		Name:        boardForm.Name,
		Description: boardForm.Description,
	}
	_, err := h.storage.CreateBoard(board)
	if err != nil {
		serverError(w, err)
		return
	}
	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}
