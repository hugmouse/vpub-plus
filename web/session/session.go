package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"
	"vpub/model"
	"vpub/storage"
)

const (
	cookieName    = "session"
	sessionMaxAge = 86400 * 30 // 30 days
)

type Manager struct {
	storage *storage.Storage
	secure  bool
}

type Session struct {
	token       string
	userID      int64
	flashErrors []string
	flashInfo   []string
}

func New(storage *storage.Storage, secure bool) *Manager {
	go storage.CleanupExpiredSessions()
	return &Manager{storage: storage, secure: secure}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *Manager) NewSession(w http.ResponseWriter, r *http.Request, userID int64) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("unable to generate session token: %w", err)
	}

	expiresAt := time.Now().Add(sessionMaxAge * time.Second)
	if err := m.storage.CreateSession(token, userID, expiresAt); err != nil {
		return nil, fmt.Errorf("unable to create session: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   sessionMaxAge,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
	})

	return &Session{token: token, userID: userID}, nil
}

func (m *Manager) GetSession(r *http.Request) (*Session, error) {
	// Probabilistic cleanup (1% of requests)
	if n, err := rand.Int(rand.Reader, big.NewInt(100)); err == nil && n.Int64() == 0 {
		go m.storage.CleanupExpiredSessions()
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, errors.New("no session cookie")
	}

	userID, err := m.storage.GetSessionUserID(cookie.Value)
	if err != nil {
		return nil, errors.New("invalid or expired session")
	}

	return &Session{token: cookie.Value, userID: userID}, nil
}

func (m *Manager) GetUser(r *http.Request) (model.User, *Session, error) {
	session, err := m.GetSession(r)
	if err != nil {
		return model.User{}, &Session{}, err
	}

	user, err := m.storage.UserByID(session.userID)
	if err != nil {
		return model.User{}, &Session{}, errors.New("user not found")
	}

	return user, session, nil
}

func (m *Manager) Delete(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}

	m.storage.DeleteSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	return nil
}

// Flash methods — buffer messages for writing to cookies in Save()

func (s *Session) FlashError(msg string) {
	s.flashErrors = append(s.flashErrors, msg)
}

func (s *Session) FlashInfo(msg string) {
	s.flashInfo = append(s.flashInfo, msg)
}

// GetFlashErrors reads flash error messages from the request cookie.
func GetFlashErrors(r *http.Request) []string {
	return readFlashCookie(r, "flash_errors")
}

// GetFlashInfo reads flash info messages from the request cookie.
func GetFlashInfo(r *http.Request) []string {
	return readFlashCookie(r, "flash_info")
}

// Save writes pending flash messages to cookies and clears consumed flash cookies.
func (s *Session) Save(r *http.Request, w http.ResponseWriter) {
	if len(s.flashErrors) > 0 {
		writeFlashCookie(w, "flash_errors", s.flashErrors)
		s.flashErrors = nil
	} else if _, err := r.Cookie("flash_errors"); err == nil {
		http.SetCookie(w, &http.Cookie{Name: "flash_errors", MaxAge: -1, Path: "/"})
	}

	if len(s.flashInfo) > 0 {
		writeFlashCookie(w, "flash_info", s.flashInfo)
		s.flashInfo = nil
	} else if _, err := r.Cookie("flash_info"); err == nil {
		http.SetCookie(w, &http.Cookie{Name: "flash_info", MaxAge: -1, Path: "/"})
	}
}

func readFlashCookie(r *http.Request, name string) []string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil
	}
	decoded, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil
	}
	var msgs []string
	if err := json.Unmarshal(decoded, &msgs); err != nil {
		return nil
	}
	return msgs
}

func writeFlashCookie(w http.ResponseWriter, name string, msgs []string) {
	encoded, err := json.Marshal(msgs)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  base64.URLEncoding.EncodeToString(encoded),
		Path:   "/",
		MaxAge: 60,
	})
}
