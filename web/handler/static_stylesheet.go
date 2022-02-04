package handler

import (
	"net/http"
	"vpub/assets"
)

func (h *Handler) showStylesheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	settings, err := h.storage.Settings()
	if err != nil {
		serverError(w, err)
		return
	}
	w.Write([]byte(assets.AssetsMap["style"] + "\n" + settings.Css))
}
