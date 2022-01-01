package model

import (
	"errors"
	"time"
)

type Post struct {
	Id        int64
	User      string
	Title     string
	Content   string
	Replies   int
	CreatedAt time.Time
	Topic     string
}

func (p Post) Validate() error {
	if len(p.Title) == 0 {
		return errors.New("title is empty")
	}
	if len(p.Content) == 0 {
		return errors.New("content is empty")
	}
	return nil
}

func (p Post) Date() string {
	return p.CreatedAt.Format("2006-01-02")
}
