package model

import (
	"errors"
	"time"
)

type Reply struct {
	Id        int64
	User      string
	Content   string
	PostId    int64
	ParentId  *int64
	PostTitle string
	Comments  int
	Thread    []Reply
	CreatedAt time.Time
	Topic     string
}

func (r Reply) Validate() error {
	if len(r.Content) == 0 {
		return errors.New("content is empty")
	}
	return nil
}

func (p Reply) Date() string {
	return p.CreatedAt.Format("2006-01-02")
}
