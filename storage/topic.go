package storage

import (
	"context"
	"vpub/model"
)

func (s *Storage) CreateTopic(userId int64, request model.TopicRequest) (int64, error) {
	var topicId int64
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return topicId, err
	}
	if err := tx.QueryRowContext(ctx, `INSERT INTO topics (is_sticky, is_locked, board_id, post_id) VALUES ($1, $2, $3, -1) RETURNING id`,
		request.IsSticky, request.IsLocked, request.BoardId).Scan(&topicId); err != nil {
		tx.Rollback()
		return topicId, err
	}
	var postId int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (subject, content, topic_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		request.Subject, request.Content, topicId, userId).Scan(&postId); err != nil {
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

func (s *Storage) UpdateTopic(id int64, request model.TopicRequest) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var postId int64
	if err := tx.QueryRowContext(ctx, `
UPDATE topics set is_locked=$1, is_sticky=$2, board_id=$3 where id=$4 returning post_id
`, request.IsLocked, request.IsSticky, request.BoardId, id).Scan(&postId); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE posts set subject=$1, content=$2 where id=$3`, request.Subject, request.Content, postId); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (s *Storage) TopicsByBoardId(boardId, page int64) ([]model.Topic, bool, error) {
	var topics []model.Topic
	settings, err := s.Settings()
	if err != nil {
		return topics, false, err
	}
	rows, err := s.db.Query(`
select topic_id,
       subject,
       content, 
       posts_count,
       updated_at,
       created_at,
       user_id,
       name,
       is_sticky
from topics_summary
where board_id=$1
offset $2
limit $3`, boardId, settings.PerPage*(page-1), settings.PerPage+1)
	if err != nil {
		return nil, false, err
	}
	for rows.Next() {
		var topic model.Topic
		err := rows.Scan(&topic.Id, &topic.Post.Subject, &topic.Post.Content, &topic.Posts, &topic.UpdatedAt, &topic.Post.CreatedAt, &topic.Post.User.Id, &topic.Post.User.Name, &topic.IsSticky)
		if err != nil {
			return topics, false, err
		}
		topics = append(topics, topic)
	}
	if len(topics) > int(settings.PerPage) {
		return topics[0:settings.PerPage], true, err
	}
	return topics, false, nil
}

func (s *Storage) TopicById(id int64) (model.Topic, error) {
	var topic model.Topic
	err := s.db.QueryRow(`
select 
       topic_id,
       posts_count,
       is_sticky,
       is_locked,
       updated_at,
       board_id,
       post_id,
       subject
from topics_summary where topic_id=$1
`, id).Scan(
		&topic.Id,
		&topic.Posts,
		&topic.IsSticky,
		&topic.IsLocked,
		&topic.UpdatedAt,
		&topic.BoardId,
		&topic.Post.Id,
		&topic.Post.Subject,
	)
	if err != nil {
		return topic, err
	}
	return topic, err
}

func (s *Storage) NewestTopicFromBoard(boardId int64) (int64, error) {
	var id int64

	query := `
        select id from topics where board_id=$1 order by updated_at desc
    `

	err := s.db.QueryRow(
		query,
		boardId,
	).Scan(
		&id,
	)

	return id, err
}
