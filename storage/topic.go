package storage

import "vpub/model"

func (s *Storage) TopicById(id int64) (model.Topic, error) {
	var topic model.Topic
	err := s.db.QueryRow(
		`SELECT id, name, description from topics WHERE id=$1`, id).Scan(
		&topic.Id,
		&topic.Name,
		&topic.Description,
	)
	return topic, err
}

func (s *Storage) Topics() ([]model.Topic, error) {
	rows, err := s.db.Query("select id, name, description from topics")
	if err != nil {
		return nil, err
	}
	var topics []model.Topic
	for rows.Next() {
		var topic model.Topic
		err := rows.Scan(&topic.Id, &topic.Name, &topic.Description)
		if err != nil {
			return topics, err
		}
		topics = append(topics, topic)
	}
	return topics, nil
}
