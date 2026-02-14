package storage

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"vpub/model"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) (string, error) {
	b := make([]rune, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[idx.Int64()]
	}
	return string(b), nil
}

func (s *Storage) CreateKey() error {
	key, err := randSeq(20)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO keys (key) VALUES ($1)`, key)
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
	defer rows.Close()
	var keys []model.Key
	for rows.Next() {
		var key model.Key
		err := rows.Scan(&key.ID, &key.Key, &key.CreatedAt)
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
		return errors.New("unable to delete key: " + err.Error())
	}

	return nil
}

func (s *Storage) KeyExists(key string) (rv bool, err error) {
	err = s.db.QueryRow(`SELECT true FROM keys WHERE key=$1 and user_id is null`, key).Scan(&rv)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return rv, nil
}
