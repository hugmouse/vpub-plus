package storage

import (
	"fmt"
	"vpub/model"
)

func (s *Storage) CreateForum(forum model.Forum) (int64, error) {
	var id int64
	err := s.db.QueryRow(`INSERT INTO forums (name, position, is_locked) VALUES ($1, $2, $3) RETURNING id`,
		forum.Name, forum.Position, forum.IsLocked).Scan(&id)
	return id, err
}

func (s *Storage) Forums() ([]model.Forum, error) {
	rows, err := s.db.Query("select id, name, position, is_locked from forums order by position")
	if err != nil {
		return nil, err
	}
	var forums []model.Forum
	for rows.Next() {
		var forum model.Forum
		err := rows.Scan(&forum.Id, &forum.Name, &forum.Position, &forum.IsLocked)
		if err != nil {
			return forums, err
		}
		forums = append(forums, forum)
	}
	return forums, nil
}

func (s *Storage) ForumById(id int64) (model.Forum, error) {
	var forum model.Forum
	err := s.db.QueryRow(
		`SELECT id, name, position, is_locked from forums WHERE id=$1`, id).Scan(
		&forum.Id,
		&forum.Name,
		&forum.Position,
		&forum.IsLocked,
	)
	return forum, err
}

func (s *Storage) UpdateForum(forum model.Forum) error {
	stmt, err := s.db.Prepare(`UPDATE forums SET name=$1, position=$2, is_locked=$3 WHERE id=$4;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(forum.Name, forum.Position, forum.IsLocked, forum.Id)
	return err
}

func (s *Storage) ForumNameExists(name string) bool {
	var result bool
	query := `SELECT true FROM forums WHERE lower(name)=lower($1) LIMIT 1`
	s.db.QueryRow(query, name).Scan(&result)
	return result
}

func (s *Storage) AnotherForumExists(id int64, name string) bool {
	var result bool
	query := `SELECT true FROM forums WHERE id != $1 AND lower(name)=lower($2) LIMIT 1`
	s.db.QueryRow(query, id, name).Scan(&result)
	return result
}
