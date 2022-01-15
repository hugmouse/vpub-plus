package model

import (
	"errors"
	"time"
)

type Post struct {
	Id        int64
	User      User
	Title     string
	Content   string
	TopicId   int64
	Replies   int
	BoardId   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Post) Validate() error {
	if len(p.Title) == 0 {
		return errors.New("title is empty")
	}
	if len(p.Title) > 120 {
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
