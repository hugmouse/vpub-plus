package storage

import (
	"math/rand"
	"vpub/model"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Storage) CreateKey() error {
	_, err := s.db.Exec(`INSERT INTO keys (key) VALUES ($1)`, randSeq(20))
	return err
}

func (s *Storage) Keys() ([]model.Key, error) {
	rows, err := s.db.Query("select key, created_at from keys where user_id is null order by created_at desc")
	if err != nil {
		return nil, err
	}
	var keys []model.Key
	for rows.Next() {
		var key model.Key
		var createdAtStr string
		err := rows.Scan(&key.Key, &createdAtStr)
		if err != nil {
			return keys, err
		}
		key.CreatedAt, err = parseCreatedAt(createdAtStr)
		if err != nil {
			return keys, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}
