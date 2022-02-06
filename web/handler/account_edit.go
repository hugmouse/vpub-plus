package handler

import (
	"net/http"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) showAccountEditPage(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserContextKey(r)

	accountForm := form.AccountForm{
		Picture: user.Picture,
		About:   user.About,
	}

	v := NewView(w, r, "account")
	v.Set("form", accountForm)
	v.Render()
}
