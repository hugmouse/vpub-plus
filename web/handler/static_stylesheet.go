package handler

import (
	"net/http"
	"vpub/assets"
	"vpub/web/handler/request"
)

func (h *Handler) showStylesheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	// Short cache: stylesheet includes dynamic admin-customizable CSS
	w.Header().Set("Cache-Control", "public, max-age=3600")

	settings := request.GetSettingsContextKey(r)

	_, err := w.Write([]byte(assets.AssetsMap["style"] + "\n" + settings.CSS))
	if err != nil {
		serverError(w, err)
	}
}
