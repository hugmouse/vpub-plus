package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"vpub/model"
)

func (s *Storage) CreateForum(request model.ForumRequest) (int64, error) {
	var id int64

	query := `
INSERT INTO forums (name, position, is_locked)
VALUES ($1, $2, $3)
RETURNING id
`
	err := s.db.QueryRow(
		query,
		request.Name,
		request.Position,
		request.IsLocked,
	).Scan(
		&id,
	)

	if err != nil {
		return id, fmt.Errorf(`store: unable to create forum %q: %v`, request.Name, err)
	}

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
		err := rows.Scan(&forum.ID, &forum.Name, &forum.Position, &forum.IsLocked)
		if err != nil {
			return forums, err
		}
		forums = append(forums, forum)
	}
	return forums, nil
}

func (s *Storage) ForumByID(id int64) (model.Forum, error) {
	var forum model.Forum
	err := s.db.QueryRow(
		`SELECT id, name, position, is_locked from forums WHERE id=$1`, id).Scan(
		&forum.ID,
		&forum.Name,
		&forum.Position,
		&forum.IsLocked,
	)
	return forum, err
}

func (s *Storage) UpdateForum(forumId int64, request model.ForumRequest) error {
	query := `
UPDATE forums 
SET name=$1, position=$2, is_locked=$3 
WHERE id=$4
`
	if _, err := s.db.Exec(query, request.Name, request.Position, request.IsLocked, forumId); err != nil {
		return errors.New("unable to update forum: " + err.Error())
	}

	return nil
}

func (s *Storage) ForumNameExists(name string) (result bool, err error) {
	query := `SELECT true FROM forums WHERE lower(name)=lower($1) LIMIT 1`
	err = s.db.QueryRow(query, name).Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return result, nil
}

func (s *Storage) AnotherForumExists(id int64, name string) (result bool, err error) {
	query := `SELECT true FROM forums WHERE id != $1 AND lower(name)=lower($2) LIMIT 1`
	err = s.db.QueryRow(query, id, name).Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return result, nil
}

func (s *Storage) RemoveForum(id int64) error {
	query := `call remove_forum($1)`
	_, err := s.db.Exec(query, id)
	return err
}
