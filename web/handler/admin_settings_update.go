package handler

import (
	"net/http"
	"vpub/model"
	"vpub/web/handler/form"
	"vpub/web/handler/request"
)

func (h *Handler) updateAdminSettings(w http.ResponseWriter, r *http.Request) {
	settingsForm := form.NewSettingsForm(r)

	if err := h.storage.UpdateSettings(model.Settings{
		Name:    settingsForm.Name,
		Css:     settingsForm.Css,
		Footer:  settingsForm.Footer,
		PerPage: settingsForm.PerPage,
	}); err != nil {
		serverError(w, err)
		return
	}

	session := request.GetSessionContextKey(r)
	session.FlashInfo("Settings updated")
	session.Save(r, w)

	http.Redirect(w, r, "/admin/settings/edit", http.StatusFound)
}
