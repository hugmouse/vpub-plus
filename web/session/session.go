package session

import (
	"github.com/gorilla/sessions"
	"net/http"
	"pboard/storage"
)

const (
	cookieName = "pboard"
	field      = "user"
)

type Session struct {
	Store   *sessions.CookieStore
	Storage *storage.Storage
}

func New(key string, storage *storage.Storage) *Session {
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteLaxMode,
	}
	return &Session{
		Store:   store,
		Storage: storage,
	}
}

func (s *Session) Save(r *http.Request, w http.ResponseWriter, name string) error {
	session, _ := s.Store.Get(r, cookieName)
	session.Values[field] = User{
		Name:            name,
		IsAuthenticated: true,
	}
	return session.Save(r, w)
}

func (s *Session) Delete(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Store.Get(r, cookieName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	return err
}

func (s *Session) Get(r *http.Request) User {
	session, err := s.Store.Get(r, cookieName)
	if err != nil {
		return User{IsAuthenticated: false}
	}
	user, ok := session.Values[field].(User)
	if !ok {
		return User{IsAuthenticated: false}
	}
	user.HasNotification = s.Storage.UserHasNotifications(user.Name)
	user.Theme = s.Storage.ThemeByUsername(user.Name)
	return user
}

type User struct {
	Name            string
	IsAuthenticated bool
	HasNotification bool
	Theme           string
}
