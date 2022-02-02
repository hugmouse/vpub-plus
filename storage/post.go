package storage

import (
	"errors"
	"vpub/model"
)

func (s *Storage) PostsByTopicId(id int64) ([]model.Post, bool, error) {
	rows, err := s.db.Query("select topic_id, post_id, subject, content, created_at, updated_at, user_id, name, picture, about from posts_full where topic_id=$1 order by created_at", id)
	if err != nil {
		return nil, false, err
	}
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.TopicId, &post.Id, &post.Subject, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.User.Id, &post.User.Name, &post.User.Picture, &post.User.About)
		if err != nil {
			return posts, false, err
		}
		posts = append(posts, post)
	}
	return posts, false, nil
}

func (s *Storage) Posts(page int64) ([]model.Post, bool, error) {
	var posts []model.Post
	settings, err := s.Settings()
	if err != nil {
		return posts, false, err
	}
	rows, err := s.db.Query(`
select
       topic_id,
       post_id,
       subject,
       created_at,
       user_id,
       name
from posts_full 
order by created_at desc
offset $1
limit $2`, settings.PerPage*(page-1), settings.PerPage+1)
	if err != nil {
		return nil, false, err
	}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.TopicId, &post.Id, &post.Subject, &post.CreatedAt, &post.User.Id, &post.User.Name)
		if err != nil {
			return posts, false, err
		}
		posts = append(posts, post)
	}
	if len(posts) > int(settings.PerPage) {
		return posts[0:settings.PerPage], true, err
	}
	return posts, false, nil
}

func (s *Storage) CreatePost(request model.PostRequest) (int64, error) {
	var id int64

	query := `
		INSERT INTO
			posts (
				   subject, 
				   content,        
				   user_id,              
				   topic_id
				   ) 
		VALUES ($1, $2, $3, $4)
		returning id
    `

	if err := s.db.QueryRow(
		query,
		request.Subject,
		request.Content,
		request.UserId,
		request.TopicId,
	).Scan(
		&id,
	); err != nil {
		return id, errors.New("unable to create post")
	}

	return id, nil
}

func (s *Storage) PostById(id int64) (model.Post, error) {
	var post model.Post
	err := s.db.QueryRow("select * from posts_full where post_id=$1", id).Scan(&post.TopicId, &post.Id, &post.Subject, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.User.Id, &post.User.Name, &post.User.Picture, &post.User.About)
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

func (s *Storage) UpdatePost(id int64, request model.PostRequest) error {
	query := `
        UPDATE posts
        SET 
            subject=$1,
            content=$2,
            updated_at=now() 
        WHERE 
              id=$3 and (user_id=$4 or (select is_admin from users where id=$4))
    `

	if _, err := s.db.Exec(
		query,
		request.Subject,
		request.Content,
		id,
		request.UserId,
	); err != nil {
		return errors.New("unable to update post")
	}

	return nil
}
