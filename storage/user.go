package storage

import (
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

func (s *Storage) CreateUser(user model.User) (int64, error) {
	var id int64
	hash, err := user.HashPassword()
	if err != nil {
		return id, err
	}
	insertUser := `INSERT INTO users (name, hash, is_admin) VALUES (lower($1), $2, $3)`
	statement, err := s.db.Prepare(insertUser)
	if err != nil {
		return id, err
	}
	_, err = statement.Exec(user.Name, hash, user.IsAdmin)
	if err != nil {
		return id, err
	}
	u, err := s.UserByName(user.Name)
	return u.Id, err
}

func (s *Storage) Users() ([]string, error) {
	rows, err := s.db.Query("select name from users")
	if err != nil {
		return nil, err
	}
	var users []string
	for rows.Next() {
		var user string
		err := rows.Scan(&user)
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
