package storage

import "vpub/model"

func (s *Storage) CreateTPost(post model.TPost) (int64, error) {
	var lid int64
	err := s.db.QueryRow(`INSERT INTO tposts (author, subject, content, topic) VALUES ($1, $2, $3, $4) RETURNING id`,
		post.Author, post.Subject, post.Content, post.Topic.Id).Scan(&lid)
	return lid, err
}
