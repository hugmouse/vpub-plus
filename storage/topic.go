package storage

import (
	"context"
	"vpub/model"
)

func (s *Storage) CreateTopic(topic model.Topic) (int64, error) {
	var topicId int64
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return topicId, err
	}
	if err := tx.QueryRowContext(ctx, `INSERT INTO topics (is_sticky, is_locked, board_id, post_id, is_sticky, is_locked) VALUES ($1, $2, $3, -1, $4, $5) RETURNING id`,
		topic.IsSticky, topic.IsLocked, topic.BoardId, topic.IsSticky, topic.IsLocked).Scan(&topicId); err != nil {
		tx.Rollback()
		return topicId, err
	}
	var postId int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (subject, content, topic_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		topic.Post.Subject, topic.Post.Content, topicId, topic.Post.User.Id).Scan(&postId); err != nil {
		tx.Rollback()
		return topicId, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE topics set post_id=$1 where id=$2`, postId, topicId); err != nil {
		tx.Rollback()
		return topicId, err
	}
	err = tx.Commit()
	return topicId, err
}

func (s *Storage) UpdateTopic(topic model.Topic) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE topics set is_locked=$1, is_sticky=$2, board_id=$3 where id=$4`, topic.IsLocked, topic.IsSticky, topic.BoardId, topic.Id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE posts set subject=$1, content=$2 where id=$3`, topic.Post.Subject, topic.Post.Content, topic.Post.Id); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (s *Storage) TopicsByBoardId(boardId int64) ([]model.Topic, bool, error) {
	rows, err := s.db.Query("select topic_id, subject, content, posts_count, updated_at, user_id, name from boardTopics where board_id=$1", boardId)
	if err != nil {
		return nil, false, err
	}
	var topics []model.Topic
	var updatedAt string
	for rows.Next() {
		var topic model.Topic
		err := rows.Scan(&topic.Id, &topic.Post.Subject, &topic.Post.Content, &topic.Replies, &updatedAt, &topic.Post.User.Id, &topic.Post.User.Name)
		topic.UpdatedAt, err = parseCreatedAt(updatedAt)
		if err != nil {
			return topics, false, err
		}
		topics = append(topics, topic)
	}
	return topics, false, nil
}

func (s *Storage) TopicById(id int64) (model.Topic, error) {
	var topic model.Topic
	var updatedAt string
	err := s.db.QueryRow(`select * from topics where id=$1`, id).Scan(&topic.Id, &topic.Replies, &topic.IsSticky, &topic.IsLocked, &updatedAt, &topic.BoardId, &topic.Post.Id)
	if err != nil {
		return topic, err
	}
	return topic, err
}
