package form

import (
	"net/http"
	"strconv"
	"strings"
)

type SettingsForm struct {
	Name                 string
	CSS                  string
	Footer               string
	PerPage              int64
	URL                  string
	Lang                 string
	SelectedRenderEngine string
	ImageProxyCacheTime  int64
	ImageProxySizeLimit  int64
	SettingsCacheTTL     int64
}

func NewSettingsForm(r *http.Request) *SettingsForm {
	perPage, _ := strconv.ParseInt(r.FormValue("per-page"), 10, 64)
	cacheTime, _ := strconv.ParseInt(r.FormValue("image-proxy-cache-time"), 10, 64)
	sizeLimit, _ := strconv.ParseInt(r.FormValue("image-proxy-size-limit"), 10, 64)
	settingsCacheTTL, _ := strconv.ParseInt(r.FormValue("settings-cache-ttl"), 10, 64)
	return &SettingsForm{
		Name:                 strings.TrimSpace(r.FormValue("name")),
		CSS:                  r.FormValue("css"),
		Footer:               r.FormValue("footer"),
		URL:                  r.FormValue("url"),
		Lang:                 r.FormValue("lang"),
		PerPage:              perPage,
		SelectedRenderEngine: r.FormValue("rendering-engine"),
		ImageProxyCacheTime:  cacheTime,
		ImageProxySizeLimit:  sizeLimit,
		SettingsCacheTTL:     settingsCacheTTL,
	}
}
