package handler

import (
	"net/http"
	"vpub/model"
	"vpub/validator"
	"vpub/web/handler/form"
)

func (h *Handler) saveAdminBoard(w http.ResponseWriter, r *http.Request) {
	boardForm := form.NewBoardForm(r)

	forums, err := h.storage.Forums()
	if err != nil {
		serverError(w, err)
		return
	}
	boardForm.Forums = forums

	v := NewView(w, r, "admin_board_create")
	v.Set("form", boardForm)

	boardRequest := model.BoardRequest{
		Name:        boardForm.Name,
		Description: boardForm.Description,
		IsLocked:    boardForm.IsLocked,
		Position:    boardForm.Position,
		ForumId:     boardForm.ForumId,
	}

	if err := validator.ValidateBoardCreation(h.storage, boardRequest); err != nil {
		v.Set("errorMessage", err.Error())
		v.Render()
		return
	}

	if _, err := h.storage.CreateBoard(boardRequest); err != nil {
		v.Set("errorMessage", "Unable to create board: "+err.Error())
		v.Render()
		return
	}

	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}
