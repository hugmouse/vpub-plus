package session

import (
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"vpub/model"
	"vpub/storage"
)

const cookieName = "vpub"

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

func (s *Session) Save(r *http.Request, w http.ResponseWriter, id int64) error {
	session, _ := s.Store.Get(r, cookieName)
	session.Values["id"] = id
	return session.Save(r, w)
}

func (s *Session) Get(r *http.Request) (model.User, error) {
	session, err := s.Store.Get(r, cookieName)
	if err != nil {
		return model.User{}, err
	}
	id, ok := session.Values["id"].(int64)
	if id == 0 || !ok {
		return model.User{}, errors.New("error extracting session")
	}
	user, err := s.Storage.UserById(id)
	if err != nil {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}
