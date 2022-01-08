package storage

import (
	"context"
	"vpub/model"
)

func (s *Storage) CreateTopic(boardId int64, post model.Post) (int64, error) {
	var topicId int64
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return topicId, err
	}
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (author, subject, content) VALUES ($1, $2, $3) RETURNING id`,
		post.User, post.Title, post.Content).Scan(&post.Id); err != nil {
		tx.Rollback()
		return topicId, err
	}
	if err := tx.QueryRowContext(ctx, `INSERT INTO topics (board_id, first_post_id) VALUES ($1, $2) RETURNING id`,
		boardId, post.Id).Scan(&topicId); err != nil {
		tx.Rollback()
		return topicId, err
	}
	err = tx.Commit()
	return topicId, err
}

func (s *Storage) TopicsByBoardId(boardId int64) ([]model.Topic, bool, error) {
	rows, err := s.db.Query("select t.id, p.subject, p.author from topics t inner join posts p on p.id = t.first_post_id where t.board_id=$1 order by t.updated_at desc;", boardId)
	if err != nil {
		return nil, false, err
	}
	var topics []model.Topic
	for rows.Next() {
		var topic model.Topic
		err := rows.Scan(&topic.Id, &topic.Subject, &topic.Author)
		if err != nil {
			return topics, false, err
		}
		topics = append(topics, topic)
	}
	return topics, false, nil
}
