package storage

import (
	"errors"
	"vpub/model"
)

func (s *Storage) BoardById(id int64) (model.Board, error) {
	var board model.Board
	err := s.db.QueryRow(`
SELECT
       b.id,
       b.name,
       description,
       b.position,
       forum_id,
       f.is_locked as forum_locked,
       b.is_locked as board_locked,
       f.name
from boards b inner join forums f on f.id = b.forum_id WHERE b.id=$1
`, id).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.Position,
		&board.Forum.Id,
		&board.Forum.IsLocked,
		&board.IsLocked,
		&board.Forum.Name,
	)
	return board, err
}

func (s *Storage) Boards() ([]model.Board, error) {
	rows, err := s.db.Query("select * from forums_summary")
	if err != nil {
		return nil, err
	}
	var boards []model.Board
	for rows.Next() {
		var board model.Board
		err := rows.Scan(&board.Id, &board.Forum.Id, &board.Forum.Name, &board.Name, &board.Description, &board.Topics, &board.Posts, &board.UpdatedAt)
		if err != nil {
			return boards, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (s *Storage) BoardsByForumId(id int64) ([]model.Board, error) {
	rows, err := s.db.Query("select * from forums_summary where forum_id=$1", id)
	if err != nil {
		return nil, err
	}
	var boards []model.Board
	for rows.Next() {
		var board model.Board
		err := rows.Scan(&board.Id, &board.Forum.Id, &board.Forum.Name, &board.Name, &board.Description, &board.Topics, &board.Posts, &board.UpdatedAt)
		if err != nil {
			return boards, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (s *Storage) CreateBoard(request model.BoardRequest) (int64, error) {
	var id int64

	query := `
INSERT INTO boards (name, description, position, forum_id, is_locked)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

	if err := s.db.QueryRow(
		query,
		request.Name,
		request.Description,
		request.Position,
		request.ForumId,
		request.IsLocked,
	).Scan(
		&id,
	); err != nil {
		return id, errors.New("unable to create board: " + err.Error())
	}

	return id, nil
}

func (s *Storage) UpdateBoard(id int64, request model.BoardRequest) error {
	query := `UPDATE boards SET name=$1, description=$2, position=$3, forum_id=$4, is_locked=$5 WHERE id=$6`

	if _, err := s.db.Exec(
		query,
		request.Name,
		request.Description,
		request.Position,
		request.ForumId,
		request.IsLocked,
		id); err != nil {
		return errors.New("unable to update board: " + err.Error())
	}

	return nil
}

func (s *Storage) BoardNameExists(name string) bool {
	var result bool
	query := `SELECT true FROM boards WHERE lower(name)=lower($1) LIMIT 1`
	s.db.QueryRow(query, name).Scan(&result)
	return result
}

func (s *Storage) AnotherBoardExists(id int64, name string) bool {
	var result bool
	query := `SELECT true FROM boards WHERE id != $1 AND lower(name)=lower($2) LIMIT 1`
	s.db.QueryRow(query, id, name).Scan(&result)
	return result
}

func (s *Storage) RemoveBoard(id int64) error {
	query := `call remove_board($1)`
	_, err := s.db.Exec(query, id)
	return err
}
