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
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (author, subject, content) VALUES ($1, $2, $3) RETURNING id`,
		post.User, post.Title, post.Content).Scan(&post.Id); err != nil {
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
	rows, err := s.db.Query("select t.id, p.subject, p.author, t.replies, t.updated_at from topics t inner join posts p on p.id = t.first_post_id where t.board_id=$1 order by t.updated_at desc;", boardId)
	if err != nil {
		return nil, false, err
	}
	var topics []model.Topic
	var updatedAtStr string
	for rows.Next() {
		var topic model.Topic
		err := rows.Scan(&topic.Id, &topic.Subject, &topic.Author, &topic.Replies, &updatedAtStr)
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
	err := s.db.QueryRow(
		`select t.id, t.board_id, t.first_post_id, t.replies, t.updated_at, p.subject, p.author from topics t left join posts p on p.id = t.first_post_id where t.id=$1;`, id).Scan(
		&topic.Id,
		&topic.BoardId,
		&topic.FirstPostId,
		&topic.Replies,
		&updatedAtStr,
		&topic.Subject,
		&topic.Author,
	)
	topic.UpdatedAt, err = parseCreatedAt(updatedAtStr)
	return topic, err
}

func (s *Storage) PostsByTopicId(id int64) ([]model.Post, bool, error) {
	rows, err := s.db.Query("select id, author, subject, content from posts where topic_id=$1;", id)
	if err != nil {
		return nil, false, err
	}
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.User, &post.Title, &post.Content)
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
	if err := tx.QueryRowContext(ctx, `INSERT INTO posts (author, subject, content, topic_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		post.User, post.Title, post.Content, post.TopicId).Scan(&post.Id); err != nil {
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
	err := s.db.QueryRow(
		`select id, author, subject, content, topic_id from posts where id=$1`, id).Scan(
		&post.Id,
		&post.User,
		&post.Title,
		&post.Content,
		&post.TopicId,
	)
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
