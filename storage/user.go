package storage

import (
	"context"
	"vpub/model"
)

const queryFindName = `SELECT id, name, hash, about, is_admin FROM users WHERE name=lower($1);`

func (s *Storage) queryUser(q string, params ...interface{}) (user model.User, err error) {
	err = s.db.QueryRow(q, params...).Scan(&user.Id, &user.Name, &user.Hash, &user.About, &user.IsAdmin)
	return
}

func (s *Storage) HasAdmin() bool {
	var rv bool
	s.db.QueryRow(`SELECT true FROM users WHERE is_admin=true limit 1`).Scan(&rv)
	return rv
}

func (s *Storage) UserHashExists(hash string) bool {
	var rv bool
	s.db.QueryRow(`SELECT true FROM users WHERE hash=$1`, hash).Scan(&rv)
	return rv
}

func (s *Storage) UserExists(name string) bool {
	var rv bool
	s.db.QueryRow(`SELECT true FROM users WHERE name=lower($1)`, name).Scan(&rv)
	return rv
}

func (s *Storage) VerifyUser(user model.User) (model.User, error) {
	u, err := s.queryUser(queryFindName, user.Name)
	if err != nil {
		return u, err
	}
	if err := user.CompareHashToPassword(u.Hash); err != nil {
		return u, err
	}
	return u, nil
}

func (s *Storage) UserByName(name string) (model.User, error) {
	return s.queryUser(queryFindName, name)
}

func (s *Storage) UserById(id int64) (model.User, error) {
	return s.queryUser(`SELECT id, name, hash, about, is_admin FROM users WHERE id=$1;`, id)
}

func (s *Storage) CreateUser(user model.User, key string) (int64, error) {
	var userId int64
	hash, err := user.HashPassword()
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return userId, err
	}
	var keyId int64
	if err := tx.QueryRowContext(ctx, `select id from keys where key=$1`, key).Scan(&keyId); err != nil {
		tx.Rollback()
		return userId, err
	}

	if err := tx.QueryRowContext(ctx, `insert into users (name, hash, is_admin, key_id) values (lower($1), $2, $3, $4) returning id`, user.Name, string(hash), user.IsAdmin, keyId).Scan(&userId); err != nil {
		tx.Rollback()
		return userId, err
	}
	if _, err := tx.ExecContext(ctx, `update keys set user_id=$1 where id=$2`, userId, keyId); err != nil {
		tx.Rollback()
		return userId, err
	}
	err = tx.Commit()
	return userId, err
}

func (s *Storage) Users() ([]model.User, error) {
	rows, err := s.db.Query("select name, hash from users")
	if err != nil {
		return nil, err
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Name, &user.Hash)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) UpdateAbout(id int64, about string) error {
	stmt, err := s.db.Prepare(`UPDATE users SET about = $1 WHERE id = $2;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(about, id)
	return err
}

func (s *Storage) UpdateUser(user model.User) error {
	stmt, err := s.db.Prepare(`UPDATE users SET name=$1, about = $2 WHERE id = $3;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Name, user.About, user.Id)
	return err
}

func (s *Storage) UpdatePassword(hash string, user model.User) error {
	newHash, err := user.HashPassword()
	stmt, err := s.db.Prepare(`UPDATE users SET hash=$1 where hash=$2;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(newHash, hash)
	return err
}
