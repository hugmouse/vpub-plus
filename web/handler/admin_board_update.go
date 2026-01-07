package handler

import (
	"fmt"
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) updateAdminBoard(w http.ResponseWriter, r *http.Request) {
	id := RouteInt64Param(r, "boardId")

	board, err := h.storage.BoardByID(id)
	if err != nil {
		serverError(w, err)
		return
	}

	boardForm := form.NewBoardForm(r)

	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}

	boardForm.Forums = forums

	v := NewView(w, r, "admin_board_create")
	v.Set("board", board)
	v.Set("form", boardForm)

	boardRequest := model.BoardRequest{
		Name:        boardForm.Name,
		Description: boardForm.Description,
		IsLocked:    boardForm.IsLocked,
		Position:    boardForm.Position,
		ForumID:     boardForm.ForumID,
	}

	if err := validator.ValidateBoardModification(h.storage, id, boardRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if err := h.storage.UpdateBoard(id, boardRequest); err != nil {
		v.Set("errorMessage", "Unable to create board: "+err.Error())
		v.Render()
		return
	}

	session := request.GetSessionContextKey(r)
	session.FlashInfo("Board updated")
	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/admin/boards/%d/edit", id), http.StatusFound)
}
