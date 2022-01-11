package storage

import (
	"context"
	"time"
	"vpub/model"
)

const postQuery = "select p.id, p.subject, p.content, p.created_at, p.topic_id, u.id, u.name from posts p left join users u on p.user_id = u.id "

func parseCreatedAt(createdAt string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", createdAt)
}

func (s *Storage) PostsByTopicId(id int64) ([]model.Post, bool, error) {
	rows, err := s.db.Query(postQuery+"where topic_id=$1;", id)
	if err != nil {
		return nil, false, err
	}
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var createdAtStr string
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &createdAtStr, &post.TopicId, &post.User.Id, &post.User.Name)
		if err != nil {
			return posts, false, err
		}
		post.CreatedAt, err = parseCreatedAt(createdAtStr)
		if err != nil {
			return posts, false, err
		}
		posts = append(posts, post)
	}
	return posts, false, nil
}

func (s *Storage) CreatePost(post model.Post) (int64, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return post.Id, err
	}
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (user_id, subject, content, topic_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		post.User.Id, post.Title, post.Content, post.TopicId).Scan(&post.Id); err != nil {
		tx.Rollback()
		return post.Id, err
	}
	var boardId int64
	if err := tx.QueryRowContext(ctx, `update topics set updated_at=datetime('now'), replies=replies + 1 where id=$1 returning board_id`, post.TopicId).Scan(&boardId); err != nil {
		tx.Rollback()
		return post.Id, err
	}
	if _, err := tx.ExecContext(ctx, `update boards set posts=posts+1, updated_at=datetime('now') where id=$1`, boardId); err != nil {
		tx.Rollback()
		return post.Id, err
	}
	err = tx.Commit()
	return post.Id, err
}

func (s *Storage) PostById(id int64) (model.Post, error) {
	var post model.Post
	var createdAtStr string
	err := s.db.QueryRow(postQuery+"where p.id=$1", id).Scan(&post.Id, &post.Title, &post.Content, &createdAtStr, &post.TopicId, &post.User.Id, &post.User.Name)
	post.CreatedAt, err = parseCreatedAt(createdAtStr)
	return post, err
}

func (s *Storage) DeletePostById(id int64) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var topicId int64
	if err := tx.QueryRowContext(ctx, `DELETE from posts WHERE id = $1 returning topic_id`, id).Scan(&topicId); err != nil {
		tx.Rollback()
		return err
	}
	var boardId int64
	if err := tx.QueryRowContext(ctx, `update topics set replies=replies - 1 where id=$1 returning board_id`, topicId).Scan(&boardId); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `update boards set posts=posts - 1 where id=$1`, boardId); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (s *Storage) UpdatePost(post model.Post) error {
	stmt, err := s.db.Prepare(`UPDATE posts SET subject = $1, content = $2 WHERE id = $3;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(post.Title, post.Content, post.Id)
	return err
}
