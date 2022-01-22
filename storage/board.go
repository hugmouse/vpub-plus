package storage

import "vpub/model"

func (s *Storage) BoardById(id int64) (model.Board, error) {
	var board model.Board
	err := s.db.QueryRow(
		`SELECT id, name, description, position, forum_id, is_locked from boards WHERE id=$1`, id).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.Position,
		&board.Forum.Id,
		&board.IsLocked,
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
		var updatedAtStr string
		err := rows.Scan(&board.Id, &board.Forum.Name, &board.Name, &board.Description, &board.Topics, &board.Posts, &updatedAtStr)
		if err != nil {
			return boards, err
		}
		board.UpdatedAt, err = parseCreatedAt(updatedAtStr)
		if err != nil {
			return boards, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (s *Storage) CreateBoard(board model.Board) (int64, error) {
	var id int64
	err := s.db.QueryRow(`INSERT INTO boards (name, description, position, forum_id, is_locked) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		board.Name, board.Description, board.Position, board.Forum.Id, board.IsLocked).Scan(&id)
	return id, err
}

func (s *Storage) UpdateBoard(board model.Board) error {
	stmt, err := s.db.Prepare(`UPDATE boards SET name=$1, description=$2, position=$3, forum_id=$4, is_locked=$5 WHERE id=$6;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(board.Name, board.Description, board.Position, board.Forum.Id, board.IsLocked, board.Id)
	return err
}
