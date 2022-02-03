package storage

import (
	"errors"
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
	rows, err := s.db.Query(`
        select 
               id,
               key,
               created_at
        from keys 
        where user_id is null
        order by created_at desc
    `)
	if err != nil {
		return nil, err
	}
	var keys []model.Key
	for rows.Next() {
		var key model.Key
		err := rows.Scan(&key.Id, &key.Key, &key.CreatedAt)
		if err != nil {
			return keys, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *Storage) DeleteKey(id int64) error {
	query := `
        delete from keys where id=$1
    `

	if _, err := s.db.Exec(
		query,
		id,
	); err != nil {
		return errors.New("unable to delete key")
	}

	return nil
}

func (s *Storage) KeyExists(key string) bool {
	var rv bool
	s.db.QueryRow(`SELECT true FROM keys WHERE key=$1 and user_id is null`, key).Scan(&rv)
	return rv
}
