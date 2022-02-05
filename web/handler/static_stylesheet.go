package handler

import (
	"net/http"
	"vpub/assets"
	"vpub/web/handler/request"
)

func (h *Handler) showStylesheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")

	settings := request.GetSettingsContextKey(r)

	w.Write([]byte(assets.AssetsMap["style"] + "\n" + settings.Css))
}
