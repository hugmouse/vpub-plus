package storage

import (
	"errors"
	"vpub/model"
)

func (s *Storage) PostsByTopicID(id int64) ([]model.Post, bool, error) {
	rows, err := s.db.Query("select topic_id, post_id, subject, content, created_at, updated_at, user_id, name, picture, about from posts_full where topic_id=$1 order by created_at", id)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.TopicID, &post.ID, &post.Subject, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.User.ID, &post.User.Name, &post.User.Picture, &post.User.About)
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
       content,
       created_at,
       updated_at,
       user_id,
       name
from posts_full 
order by created_at desc
offset $1
limit $2`, settings.PerPage*(page-1), settings.PerPage+1)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.TopicID,
			&post.ID,
			&post.Subject,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.User.ID,
			&post.User.Name,
		)
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

func (s *Storage) PostsByUserID(id, page int64) ([]model.Post, bool, error) {
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
       content,
       created_at,
       updated_at,
       user_id,
       name
from posts_full 
where user_id=$1
order by created_at desc
offset $2
limit $3`, id, settings.PerPage*(page-1), settings.PerPage+1)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.TopicID,
			&post.ID,
			&post.Subject,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.User.ID,
			&post.User.Name,
		)
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

func (s *Storage) CreatePost(userId, topicId int64, request model.PostRequest) (int64, error) {
	var id int64

	query := `
		INSERT INTO	posts 
		    (subject, content, user_id, topic_id) 
		VALUES ($1, $2, $3, $4)
		returning id
    `

	if err := s.db.QueryRow(
		query,
		request.Subject,
		request.Content,
		userId,
		topicId,
	).Scan(
		&id,
	); err != nil {
		return id, errors.New("unable to create post: " + err.Error())
	}

	return id, nil
}

func (s *Storage) PostByID(id int64) (model.Post, error) {
	var post model.Post
	err := s.db.QueryRow("select * from posts_full where post_id=$1 limit 1", id).Scan(&post.TopicID, &post.ID, &post.Subject, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.User.ID, &post.User.Name, &post.User.Picture, &post.User.About)
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
	defer stmt.Close()
	_, err = stmt.Exec(post.ID, post.User.ID)
	return err
}

func (s *Storage) UpdatePost(id, userId int64, request model.PostRequest) error {
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
		userId,
	); err != nil {
		return errors.New("unable to update post: " + err.Error())
	}

	return nil
}

func (s *Storage) NewestPostFromTopic(topicId int64) (int64, error) {
	var id int64

	query := `
        select id from posts where topic_id=$1 order by created_at desc limit 1
    `

	err := s.db.QueryRow(
		query,
		topicId,
	).Scan(
		&id,
	)

	return id, err
}
