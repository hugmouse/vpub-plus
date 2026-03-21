package storage

import (
	"log"
	"time"
)

func (s *Storage) CreateSession(token string, userID int64, expiresAt time.Time) error {
	_, err := s.db.Exec(
		`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, expiresAt,
	)
	return err
}

func (s *Storage) GetSessionUserID(token string) (int64, error) {
	var userID int64
	err := s.db.QueryRow(
		`SELECT user_id FROM sessions WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&userID)
	return userID, err
}

func (s *Storage) DeleteSession(token string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func (s *Storage) CleanupExpiredSessions() {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token IN (SELECT token FROM sessions WHERE expires_at < NOW() LIMIT 1000)`)
	if err != nil {
		log.Println("[session cleanup] error:", err)
	}
}
