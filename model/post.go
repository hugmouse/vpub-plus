package model

import (
	"time"
)

type Post struct {
	Id        int64
	User      User
	Subject   string
	Content   string
	TopicId   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PostRequest represents the request to create or update a post.
type PostRequest struct {
	Subject string
	Content string
}

func (p Post) Date() string {
	return p.CreatedAt.Format("2006-01-02")
}

func (p Post) DateUpdated() string {
	return p.UpdatedAt.Format("2006-01-02")
}
