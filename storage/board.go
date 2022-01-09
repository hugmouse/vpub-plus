package storage

import "vpub/model"

func (s *Storage) BoardById(id int64) (model.Board, error) {
	var board model.Board
	err := s.db.QueryRow(
		`SELECT id, name, description, topics, posts from boards WHERE id=$1`, id).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.Topics,
		&board.Posts,
	)
	return board, err
}

func (s *Storage) Boards() ([]model.Board, error) {
	//	var createdAtStr string
	//	var updatedAtStr string
	//	err := rows.Scan(&post.Id, &post.User, &post.Title, &createdAtStr, &updatedAtStr, &post.Topic)
	//	post.CreatedAt, err = parseCreatedAt(createdAtStr)
	//	post.UpdatedAt, err = parseCreatedAt(updatedAtStr)
	//	if err != nil {
	//		return post, err
	//	}
	//	return post, nil
	rows, err := s.db.Query("select id, name, description, topics, posts, updated_at from boards")
	if err != nil {
		return nil, err
	}
	var boards []model.Board
	for rows.Next() {
		var board model.Board
		var updatedAtStr string
		err := rows.Scan(&board.Id, &board.Name, &board.Description, &board.Topics, &board.Posts, &updatedAtStr)
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
