package model

import (
	"errors"
	"time"
)

type TPost struct {
	Id        int64
	Author    string
	Subject   string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Topic     Board
}

func (p TPost) Validate() error {
	if len(p.Subject) == 0 {
		return errors.New("subject is empty")
	}
	if len(p.Content) == 0 {
		return errors.New("content is empty")
	}
	return nil
}

func (p TPost) DateCreated() string {
	return p.CreatedAt.Format("2006-01-02")
}

func (p TPost) DateUpdated() string {
	return p.UpdatedAt.Format("2006-01-02")
}
