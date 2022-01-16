package storage

import "vpub/model"

func (s *Storage) CreateForum(forum model.Forum) (int64, error) {
	var id int64
	err := s.db.QueryRow(`INSERT INTO forums (name, position) VALUES ($1, $2) RETURNING id`,
		forum.Name, forum.Position).Scan(&id)
	return id, err
}

func (s *Storage) Forums() ([]model.Forum, error) {
	rows, err := s.db.Query("select id, name, position from forums order by position")
	if err != nil {
		return nil, err
	}
	var forums []model.Forum
	for rows.Next() {
		var forum model.Forum
		err := rows.Scan(&forum.Id, &forum.Name, &forum.Position)
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
		`SELECT id, name, position from forums WHERE id=$1`, id).Scan(
		&forum.Id,
		&forum.Name,
		&forum.Position,
	)
	return forum, err
}

func (s *Storage) UpdateForum(forum model.Forum) error {
	stmt, err := s.db.Prepare(`UPDATE forums SET name=$1, position=$2 WHERE id=$3;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(forum.Name, forum.Position, forum.Id)
	return err
}
