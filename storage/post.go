package storage

import (
	"vpub/model"
)

func (s *Storage) PostsByTopicId(id int64) ([]model.Post, bool, error) {
	rows, err := s.db.Query("select topic_id, post_id, subject, content, created_at, user_id, name, picture from postUsers where topic_id=$1", id)
	if err != nil {
		return nil, false, err
	}
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.TopicId, &post.Id, &post.Subject, &post.Content, &post.CreatedAt, &post.User.Id, &post.User.Name, &post.User.Picture)
		if err != nil {
			return posts, false, err
		}
		posts = append(posts, post)
	}
	return posts, false, nil
}

func (s *Storage) CreatePost(post model.Post) (int64, error) {
	var id int64
	err := s.db.QueryRow(`INSERT INTO posts (subject, content, user_id, topic_id) VALUES ($1, $2, $3, $4) returning id`, post.Subject, post.Content, post.User.Id, post.TopicId).Scan(&id)
	return id, err
}

func (s *Storage) PostById(id int64) (model.Post, error) {
	var post model.Post
	err := s.db.QueryRow("select * from postUsers where post_id=$1", id).Scan(&post.TopicId, &post.Id, &post.Subject, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.User.Id, &post.User.Name, &post.User.Picture)
	if err != nil {
		return post, err
	}
	return post, err
}

func (s *Storage) DeletePost(post model.Post) error {
	stmt, err := s.db.Prepare(`delete from posts where id=$1 and (user_id = $2 or (select is_admin from users where id=$2))`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(post.Id, post.User.Id)
	return err
}

func (s *Storage) UpdatePost(post model.Post) error {
	stmt, err := s.db.Prepare(`UPDATE posts SET subject=$1, content=$2, updated_at=datetime('now'), is_sticky=$3, is_locked=$4, board_id=$5 WHERE id=$6 and (user_id=$7 or (select is_admin from users where id=$7));`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(post.Subject, post.Content, post.Content, post.Content, post.Content, post.Id, post.User.Id)
	return err
}
