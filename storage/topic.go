package storage

import (
	"context"
	"vpub/model"
)

func (s *Storage) CreateTopic(topic model.Topic) (int64, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return topic.Id, err
	}
	post := topic.FirstPost
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (author, subject, content) VALUES ($1, $2, $3) RETURNING id`,
		post.User, post.Title, post.Content).Scan(&post.Id); err != nil {
		tx.Rollback()
		return topic.Id, err
	}
	if err := tx.QueryRowContext(ctx, `INSERT INTO topics (board_id, first_post_id) VALUES ($1, $2) RETURNING id`,
		topic.BoardId, post.Id).Scan(&topic.Id); err != nil {
		tx.Rollback()
		return topic.Id, err
	}
	err = tx.Commit()
	return topic.Id, err
}
