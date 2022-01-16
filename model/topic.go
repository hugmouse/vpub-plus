package model

import "time"

type Topic struct {
	Id        int64
	BoardId   int64
	Subject   string
	User      User
	IsSticky  bool
	Replies   int64
	UpdatedAt time.Time
	CreatedAt time.Time
}
