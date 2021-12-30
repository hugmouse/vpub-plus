package storage

import (
	"database/sql"
	"pboard/model"
	"strconv"
	"strings"
	"time"
)

type postQueryBuilder struct {
	where  string
	limit  string
	offset string
}

func (p postQueryBuilder) build() string {
	query := []string{`SELECT id, author, title, created_at from posts`}
	if p.where != "" {
		query = append(query, `WHERE`, p.where)
	}
	query = append(query, `ORDER BY created_at desc`)
	if p.limit != "" {
		query = append(query, `LIMIT`, p.limit)
	}
	if p.offset != "" {
		query = append(query, `OFFSET`, p.offset)
	}
	return strings.Join(query, " ")
}

func (s *Storage) populatePost(rows *sql.Rows) (model.Post, error) {
	var post model.Post
	var createdAtStr string
	err := rows.Scan(&post.Id, &post.User, &post.Title, &createdAtStr)
	post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (s *Storage) CreatePost(post model.Post) (int64, error) {
	var lid int64
	err := s.db.QueryRow(`INSERT INTO posts (author, title, content) VALUES ($1, $2, $3) RETURNING id`,
		post.User, post.Title, post.Content).Scan(&lid)
	return lid, err
}

func (s *Storage) PostById(id int64) (model.Post, error) {
	var post model.Post
	err := s.db.QueryRow(
		`SELECT id, author, title, content from posts WHERE id=$1`, id).Scan(
		&post.Id,
		&post.User,
		&post.Title,
		&post.Content,
	)
	return post, err
}

func (s *Storage) PostsByUsername(user string, perPage int, page int64) ([]model.Post, bool, error) {
	rows, err := s.db.Query(postQueryBuilder{
		where:  `author = $1`,
		limit:  strconv.Itoa(perPage + 1),
		offset: `$2`,
	}.build(), user, page*int64(perPage))
	if err != nil {
		return nil, false, err
	}
	var posts []model.Post
	for rows.Next() {
		post, err := s.populatePost(rows)
		if err != nil {
			return posts, false, err
		}
		posts = append(posts, post)
	}
	if len(posts) > perPage {
		return posts[0:perPage], true, err
	}
	return posts, false, err
}

func (s *Storage) Posts(page int64, perPage int) ([]model.Post, bool, error) {
	rows, err := s.db.Query(postQueryBuilder{
		limit:  strconv.Itoa(perPage + 1),
		offset: `$1`,
	}.build(), page*int64(perPage))
	if err != nil {
		return nil, false, err
	}
	var posts []model.Post
	for rows.Next() {
		post, err := s.populatePost(rows)
		if err != nil {
			return posts, false, err
		}
		posts = append(posts, post)
	}
	if len(posts) > perPage {
		return posts[0:perPage], true, err
	}
	return posts, false, err
}

func (s *Storage) UpdatePost(post model.Post) error {
	stmt, err := s.db.Prepare(`UPDATE posts SET title = $1, content = $2 WHERE id = $3 and author = $4;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(post.Title, post.Content, post.Id, post.User)
	return err
}

func (s *Storage) DeletePost(id int64, author string) error {
	stmt, err := s.db.Prepare(`DELETE from posts WHERE id = $1 and author = $2;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, author)
	return err
}
