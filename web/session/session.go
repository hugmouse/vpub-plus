package session

import (
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"vpub/model"
	"vpub/storage"
)

const cookieName = "status"

type Manager struct {
	Store   *sessions.CookieStore
	Storage *storage.Storage
}

type Session struct {
	session *sessions.Session
}

func (s *Manager) NewSession(w http.ResponseWriter, r *http.Request, userID int64) (*Session, error) {
	sess, err := s.Store.Get(r, cookieName)
	if err != nil {
		return nil, fmt.Errorf("unable to get session: %w", err)
	}
	sess.Values["id"] = userID

	session := &Session{session: sess}
	if err := session.session.Save(r, w); err != nil {
		return nil, fmt.Errorf("unable to save session: %w", err)
	}
	return session, nil
}

func (s *Session) FlashError(msg string) {
	s.session.AddFlash(msg, "errors")
}

func (s *Session) FlashInfo(msg string) {
	s.session.AddFlash(msg, "info")
}

func (s *Session) Save(r *http.Request, w http.ResponseWriter) {
	if err := s.session.Save(r, w); err != nil {
		fmt.Println("error saving session")
	}
}

func (s *Session) SetUserId(id int64) {
	s.session.Values["id"] = id
}

func (s *Session) GetFlashErrors() []string {
	var errorsArray []string
	if msgs := s.session.Flashes("errors"); len(msgs) > 0 {
		for _, m := range msgs {
			errorsArray = append(errorsArray, m.(string))
		}
	}
	return errorsArray
}

func (s *Session) GetFlashInfo() []string {
	var info []string
	if msgs := s.session.Flashes("info"); len(msgs) > 0 {
		for _, m := range msgs {
			info = append(info, m.(string))
		}
	}
	return info
}

func New(key string, storage *storage.Storage) *Manager {
	store := sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   86400 * 30,
		Path:     "/",
	}
	return &Manager{
		Store:   store,
		Storage: storage,
	}
}

func (s *Manager) Delete(w http.ResponseWriter, r *http.Request) error {
	session, err := s.GetSession(r)
	if err != nil {
		return err
	}
	session.session.Options.MaxAge = -1
	session.Save(r, w)
	return err
}

func (s *Manager) GetSession(r *http.Request) (*Session, error) {
	sess, err := s.Store.Get(r, cookieName)
	if err != nil {
		return nil, err
	}
	return &Session{session: sess}, nil
}

//func (s *Manager) Save(r *http.Request, w http.ResponseWriter, id int64) error {
//	session, _ := s.GetSession(r)
//	session.Values["id"] = id
//	return session.Save(r, w)
//}

// GetUser Returns an error if the user doesn't exist
func (s *Manager) GetUser(r *http.Request) (model.User, *Session, error) {
	session, err := s.GetSession(r)
	if err != nil {
		return model.User{}, &Session{}, err
	}
	id, ok := session.session.Values["id"].(int64)
	if id == 0 || !ok {
		return model.User{}, &Session{}, errors.New("error extracting session")
	}
	user, err := s.Storage.UserById(id)
	if err != nil {
		return model.User{}, &Session{}, errors.New("user not found")
	}
	return user, session, nil
}
