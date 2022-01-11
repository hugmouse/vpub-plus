package model

import "time"

type Topic struct {
	Id          int64
	BoardId     int64
	Subject     string
	User        User
	FirstPostId int64
	Replies     int64
	UpdatedAt   time.Time
}
