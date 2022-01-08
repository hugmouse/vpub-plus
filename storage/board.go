package storage

import "vpub/model"

func (s *Storage) BoardById(id int64) (model.Board, error) {
	var board model.Board
	err := s.db.QueryRow(
		`SELECT id, name, description from boards WHERE id=$1`, id).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
	)
	return board, err
}

func (s *Storage) Boards() ([]model.Board, error) {
	rows, err := s.db.Query("select id, name, description from boards")
	if err != nil {
		return nil, err
	}
	var boards []model.Board
	for rows.Next() {
		var board model.Board
		err := rows.Scan(&board.Id, &board.Name, &board.Description)
		if err != nil {
			return boards, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}
