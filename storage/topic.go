package storage

import (
	"vpub/model"
)

func (s *Storage) CreateTopic(post model.Post) (int64, error) {
	var id int64
	err := s.db.QueryRow(`INSERT INTO posts (subject, content, board_id, is_sticky, user_id) VALUES ($1, $2, $3, $4, $5) returning id`, post.Title, post.Content, post.BoardId, post.IsSticky, post.User.Id).Scan(&id)
	return id, err
}

func (s *Storage) TopicsByBoardId(boardId int64) ([]model.Topic, bool, error) {
	rows, err := s.db.Query("select * from topics where board_id=$1", boardId)
	if err != nil {
		return nil, false, err
	}
	var topics []model.Topic
	var updatedAt string
	var createdAt string
	for rows.Next() {
		var topic model.Topic
		err := rows.Scan(&topic.Id, &topic.User.Id, &topic.User.Name, &topic.BoardId, &topic.Subject, &topic.Replies, &createdAt, &updatedAt, &topic.IsSticky)
		topic.CreatedAt, err = parseCreatedAt(createdAt)
		if err != nil {
			return topics, false, err
		}
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
	var createdAt string
	err := s.db.QueryRow(`select * from topics where id=$1`, id).Scan(&topic.Id, &topic.User.Id, &topic.User.Name, &topic.BoardId, &topic.Subject, &topic.Replies, &createdAt, &updatedAt, &topic.IsSticky)
	topic.CreatedAt, err = parseCreatedAt(createdAt)
	if err != nil {
		return topic, err
	}
	topic.UpdatedAt, err = parseCreatedAt(updatedAt)
	if err != nil {
		return topic, err
	}
	return topic, err
}
