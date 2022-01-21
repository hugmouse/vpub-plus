package model

import "time"

type Topic struct {
	Id        int64
	BoardId   int64
	IsSticky  bool
	IsLocked  bool
	Replies   int64
	UpdatedAt time.Time
	Post      Post
}
