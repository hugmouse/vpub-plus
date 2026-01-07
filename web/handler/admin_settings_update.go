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
		Name:                settingsForm.Name,
		CSS:                 settingsForm.CSS,
		Footer:              settingsForm.Footer,
		PerPage:             settingsForm.PerPage,
		URL:                 settingsForm.URL,
		Lang:                settingsForm.Lang,
		ImageProxyCacheTime: settingsForm.ImageProxyCacheTime,
		ImageProxySizeLimit: settingsForm.ImageProxySizeLimit,
	}); err != nil {
		serverError(w, err)
		return
	}

	engine, err := h.renderRegistry.Get(settingsForm.SelectedRenderEngine)
	if err != nil {
		serverError(w, err)
		return
	}
	h.currentRenderEngine = &engine

	session := request.GetSessionContextKey(r)
	session.FlashInfo("Settings updated")
	session.Save(r, w)

	http.Redirect(w, r, "/admin/settings/edit", http.StatusFound)
}
