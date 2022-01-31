package request

import (
	"net/http"
	"vpub/model"
	"vpub/web/session"
)

const (
	UserKey = iota
	SessionKey
	SettingsKey
)

func GetSettingsContextKey(r *http.Request) model.Settings {
	settings, ok := r.Context().Value(SettingsKey).(model.Settings)
	if !ok {
		return model.Settings{}
	}
	return settings
}

func GetSessionContextKey(r *http.Request) *session.Session {
	session, ok := r.Context().Value(SessionKey).(*session.Session)
	if !ok {
		return nil
	}
	return session
}

func GetUserContextKey(r *http.Request) model.User {
	user, ok := r.Context().Value(UserKey).(model.User)
	if !ok {
		return model.User{}
	}
	return user
}
