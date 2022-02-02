package model

import (
	"errors"
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
	UserId  int64
	Subject string
	Content string
	TopicId int64
}

func (p Post) Validate() error {
	if len(p.Subject) == 0 {
		return errors.New("title is empty")
	}
	if len(p.Subject) > 120 {
		return errors.New("title has a max length of 120 characters")
	}
	if len(p.Content) == 0 {
		return errors.New("content is empty")
	}
	return nil
}

func (p Post) Date() string {
	return p.CreatedAt.Format("2006-01-02")
}

func (p Post) DateUpdated() string {
	return p.UpdatedAt.Format("2006-01-02")
}
