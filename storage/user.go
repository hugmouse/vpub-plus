package storage

import (
	"pboard/model"
)

const queryFindName = `SELECT name, hash, about FROM users WHERE name=lower($1);`

func (s *Storage) queryUser(q string, params ...interface{}) (user model.User, err error) {
	err = s.db.QueryRow(q, params...).Scan(&user.Name, &user.Hash, &user.About)
	return
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

func (s *Storage) CreateUser(user model.User) error {
	hash, err := user.HashPassword()
	if err != nil {
		return err
	}
	insertUser := `INSERT INTO users (name, hash) VALUES (lower($1), $2)`
	statement, err := s.db.Prepare(insertUser)
	if err != nil {
		return err
	}
	_, err = statement.Exec(user.Name, hash)
	return err
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

func (s *Storage) RandomUsers(n int) ([]string, error) {
	rows, err := s.db.Query("select  author from posts group by author order by random() limit $1", n)
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

func (s *Storage) UpdateTheme(username, name string) error {
	stmt, err := s.db.Prepare(`UPDATE users SET theme = $1 WHERE name = $2;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, username)
	return err
}

func (s *Storage) UpdateAbout(username, about string) error {
	stmt, err := s.db.Prepare(`UPDATE users SET about = $1 WHERE name = $2;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(about, username)
	return err
}

func (s *Storage) ThemeByUsername(name string) string {
	var rv string
	s.db.QueryRow(`SELECT theme FROM users WHERE name=$1`, name).Scan(&rv)
	return rv
}
