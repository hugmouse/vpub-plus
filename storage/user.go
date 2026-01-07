package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"vpub/model"
)

const queryFindName = `SELECT id, name, hash, about, is_admin, picture FROM users WHERE name=lower($1);`

type ErrUserNotFound struct{}

func (u ErrUserNotFound) Error() string {
	return "user not found"
}

type ErrWrongPassword struct{}

func (u ErrWrongPassword) Error() string {
	return "wrong password"
}

type ErrUserExists struct{}

func (u ErrUserExists) Error() string {
	return "user already exists"
}

func (s *Storage) queryUser(q string, params ...interface{}) (user model.User, err error) {
	err = s.db.QueryRow(q, params...).Scan(&user.Id, &user.Name, &user.Hash, &user.About, &user.IsAdmin, &user.Picture)
	return
}

func (s *Storage) HasAdmin() (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT true FROM users WHERE is_admin=true LIMIT 1`).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (s *Storage) UserHashExists(hash string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT true FROM users WHERE hash=$1`, hash).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (s *Storage) UserExists(name string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT true FROM users WHERE name=LOWER($1)`, name).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (s *Storage) VerifyUser(user model.User) (model.User, error) {
	u, err := s.queryUser(queryFindName, user.Name)
	if err != nil {
		return u, ErrUserNotFound{}
	}
	if err := user.CompareHashToPassword(u.Hash); err != nil {
		return u, ErrWrongPassword{}
	}
	return u, nil
}

func (s *Storage) UserByName(name string) (model.User, error) {
	return s.queryUser(queryFindName, name)
}

func (s *Storage) UserById(id int64) (model.User, error) {
	return s.queryUser(`SELECT id, name, hash, about, is_admin, picture FROM users WHERE id=$1;`, id)
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}

func (s *Storage) CreateUser(key string, request model.UserCreationRequest) (int64, error) {
	var userId int64
	hash, err := hashPassword(request.Password)
	if err != nil {
		return userId, fmt.Errorf("failed to hash password: %w", err)
	}
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return userId, err
	}
	var keyId int64
	if err := tx.QueryRowContext(ctx, `select id from keys where key=$1 and user_id is null`, key).Scan(&keyId); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return userId, errors.Join(err, fmt.Errorf("rollback in CreateUser failed: %w", rbErr))
		}
		return userId, err
	}
	var exists bool
	if err := tx.QueryRowContext(ctx, `select exists(select 1 from users where name=$1)`, request.Name).Scan(&exists); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return userId, errors.Join(err, fmt.Errorf("rollback in CreateUser failed: %w", rbErr))
		}
		return userId, err
	}
	if exists {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return userId, errors.Join(err, fmt.Errorf("rollback in CreateUser failed: %w", rbErr))
		}
		return userId, ErrUserExists{}
	}
	if err := tx.QueryRowContext(ctx, `insert into users (name, hash, is_admin) values (lower($1), $2, $3) returning id`, request.Name, string(hash), request.IsAdmin).Scan(&userId); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return userId, errors.Join(err, fmt.Errorf("rollback in CreateUser failed: %w", rbErr))
		}
		return userId, err
	}
	if _, err := tx.ExecContext(ctx, `update keys set user_id=$1 where id=$2`, userId, keyId); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return userId, errors.Join(err, fmt.Errorf("rollback in CreateUser failed: %w", rbErr))
		}
		return userId, err
	}
	err = tx.Commit()
	return userId, err
}

func (s *Storage) Users() ([]model.User, error) {
	rows, err := s.db.Query("select id, name, hash from users")
	if err != nil {
		return nil, err
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Hash)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) UpdateUser(user model.User) error {
	stmt, err := s.db.Prepare(`UPDATE users SET name=$1, about=$2, picture=$3 WHERE id = $4;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Name, user.About, user.Picture, user.Id)
	return err
}

func (s *Storage) UpdatePassword(hash string, user model.User) error {
	newHash, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}
	stmt, err := s.db.Prepare(`UPDATE users SET hash=$1 where hash=$2;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(newHash, hash)
	return err
}

func (s *Storage) RemoveUser(id int64) error {
	query := `call remove_user($1)`
	_, err := s.db.Exec(query, id)
	return err
}
