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
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (user_id, subject, content) VALUES ($1, $2, $3) RETURNING id`,
		post.User.Id, post.Title, post.Content).Scan(&post.Id); err != nil {
		tx.Rollback()
		return topicId, err
	}
	if err := tx.QueryRowContext(ctx, `INSERT INTO topics (board_id, first_post_id) VALUES ($1, $2) RETURNING id`,
		boardId, post.Id).Scan(&topicId); err != nil {
		tx.Rollback()
		return topicId, err
	}
	if _, err := tx.ExecContext(ctx, `update posts set topic_id=$1 where id=$2`, topicId, post.Id); err != nil {
		tx.Rollback()
		return topicId, err
	}
	if _, err := tx.ExecContext(ctx, `update boards set topics=topics+1, posts=posts+1, updated_at=datetime('now') where id=$1`, boardId); err != nil {
		tx.Rollback()
		return topicId, err
	}
	err = tx.Commit()
	return topicId, err
}

func (s *Storage) TopicsByBoardId(boardId int64) ([]model.Topic, bool, error) {
	rows, err := s.db.Query("select t.id, p.subject, t.replies, t.updated_at, u.id, u.name from topics t inner join posts p on p.id = t.first_post_id inner join users u on u.id = p.user_id where t.board_id=$1 order by t.updated_at desc;", boardId)
	if err != nil {
		return nil, false, err
	}
	var topics []model.Topic
	var updatedAtStr string
	for rows.Next() {
		var topic model.Topic
		err := rows.Scan(&topic.Id, &topic.Subject, &topic.Replies, &updatedAtStr, &topic.User.Id, &topic.User.Name)
		topic.UpdatedAt, err = parseCreatedAt(updatedAtStr)
		if err != nil {
			return topics, false, err
		}
		topic.UpdatedAt, err = parseCreatedAt(updatedAtStr)
		if err != nil {
			return topics, false, err
		}
		topics = append(topics, topic)
	}
	return topics, false, nil
}

func (s *Storage) TopicById(id int64) (model.Topic, error) {
	var topic model.Topic
	var updatedAtStr string
	err := s.db.QueryRow(`select t.id, p.subject, t.replies, t.updated_at, t.board_id, t.first_post_id, u.id, u.name from topics t inner join posts p on p.id = t.first_post_id inner join users u on u.id = p.user_id where t.id=$1;`, id).Scan(&topic.Id, &topic.Subject, &topic.Replies, &updatedAtStr, &topic.BoardId, &topic.FirstPostId, &topic.User.Id, &topic.User.Name)
	topic.UpdatedAt, err = parseCreatedAt(updatedAtStr)
	return topic, err
}
