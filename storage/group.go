package storage

import (
	"database/sql"
	"errors"
	"vpub/model"
)

func (s *Storage) CreateGroup(req model.GroupRequest) (int64, error) {
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO groups (name) VALUES ($1) RETURNING id`,
		req.Name,
	).Scan(&id)
	return id, err
}

func (s *Storage) Groups() ([]model.Group, error) {
	rows, err := s.db.Query(`
		SELECT g.id, g.name, COUNT(gm.user_id) AS member_count
		FROM groups g
		LEFT JOIN group_members gm ON gm.group_id = g.id
		GROUP BY g.id, g.name
		ORDER BY g.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []model.Group
	for rows.Next() {
		var g model.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.MemberCount); err != nil {
			return groups, err
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return groups, err
	}
	return groups, nil
}

func (s *Storage) GroupByID(id int64) (model.Group, error) {
	var g model.Group
	err := s.db.QueryRow(
		`SELECT id, name FROM groups WHERE id=$1`, id,
	).Scan(&g.ID, &g.Name)
	return g, err
}

func (s *Storage) UpdateGroup(id int64, req model.GroupRequest) error {
	_, err := s.db.Exec(`UPDATE groups SET name=$1 WHERE id=$2`, req.Name, id)
	return err
}

func (s *Storage) RemoveGroup(id int64) error {
	_, err := s.db.Exec(`DELETE FROM groups WHERE id=$1`, id)
	return err
}

func (s *Storage) GroupNameExists(name string) (result bool, err error) {
	err = s.db.QueryRow(
		`SELECT true FROM groups WHERE lower(name)=lower($1) LIMIT 1`, name,
	).Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return result, err
}

func (s *Storage) AnotherGroupExists(id int64, name string) (result bool, err error) {
	err = s.db.QueryRow(
		`SELECT true FROM groups WHERE id != $1 AND lower(name)=lower($2) LIMIT 1`, id, name,
	).Scan(&result)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return result, err
}

func (s *Storage) GroupMembers(groupID int64) ([]model.User, error) {
	rows, err := s.db.Query(`
		SELECT u.id, u.name
		FROM users u
		JOIN group_members gm ON gm.user_id = u.id
		WHERE gm.group_id = $1
		ORDER BY u.name
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return users, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (s *Storage) AddGroupMember(groupID, userID int64) error {
	_, err := s.db.Exec(
		`INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		groupID, userID,
	)
	return err
}

func (s *Storage) RemoveGroupMember(groupID, userID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM group_members WHERE group_id=$1 AND user_id=$2`,
		groupID, userID,
	)
	return err
}

func (s *Storage) ForumsWithGroupCount(groupID int64) (int64, error) {
	var count int64
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM forums WHERE group_id=$1`, groupID,
	).Scan(&count)
	return count, err
}
