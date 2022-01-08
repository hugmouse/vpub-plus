package storage

import "vpub/model"

func (s *Storage) CreateTPost(post model.TPost) (int64, error) {
	var lid int64
	err := s.db.QueryRow(`INSERT INTO tposts (author, subject, content, topic) VALUES ($1, $2, $3, $4) RETURNING id`,
		post.Author, post.Subject, post.Content, post.Topic.Id).Scan(&lid)
	return lid, err
}

func (s *Storage) ThreadsByTopicId(topicId int64) ([]model.TPost, bool, error) {
	rows, err := s.db.Query("select id, author, subject, content from tposts")
	if err != nil {
		return nil, false, err
	}
	var threads []model.TPost
	for rows.Next() {
		var thread model.TPost
		err := rows.Scan(&thread.Id, &thread.Author, &thread.Subject, &thread.Content)
		if err != nil {
			return threads, false, err
		}
		threads = append(threads, thread)
	}
	return threads, false, nil
}
