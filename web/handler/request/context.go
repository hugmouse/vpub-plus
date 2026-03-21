package request

import (
	"net/http"
	"vpub/model"
	"vpub/web/session"
)

type contextKey string

const (
	UserKey      contextKey = "user"
	SessionKey   contextKey = "session"
	SettingsKey  contextKey = "settings"
	CSRFTokenKey contextKey = "csrfToken"
)

func GetSettingsContextKey(r *http.Request) model.Settings {
	settings, ok := r.Context().Value(SettingsKey).(model.Settings)
	if !ok {
		return model.Settings{}
	}
	return settings
}

func GetSessionContextKey(r *http.Request) *session.Session {
	sessionFromContext, ok := r.Context().Value(SessionKey).(*session.Session)
	if !ok {
		return nil
	}
	return sessionFromContext
}

func GetUserContextKey(r *http.Request) model.User {
	user, ok := r.Context().Value(UserKey).(model.User)
	if !ok {
		return model.User{}
	}
	return user
}

func GetCSRFTokenContextKey(r *http.Request) string {
	token, _ := r.Context().Value(CSRFTokenKey).(string)
	return token
}
