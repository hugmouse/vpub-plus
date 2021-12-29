package session

import (
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"pboard/storage"
)

const cookieName = "pboard"

type Session struct {
	Store   *sessions.CookieStore
	Storage *storage.Storage
}

func New(key string, storage *storage.Storage) *Session {
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
	return &Session{
		Store:   store,
		Storage: storage,
	}
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

func (s *Session) Save(r *http.Request, w http.ResponseWriter, name string) error {
	session, _ := s.Store.Get(r, cookieName)
	session.Values["name"] = name
	return session.Save(r, w)
}

func (s *Session) Get(r *http.Request) (string, error) {
	session, err := s.Store.Get(r, cookieName)
	if err != nil {
		return "", err
	}
	name, ok := session.Values["name"].(string)
	if name == "" || !ok {
		return "", errors.New("error extracting session")
	}
	if ok = s.Storage.UserExists(name); !ok {
		return "", errors.New("username not found")
	}
	return name, nil
}
