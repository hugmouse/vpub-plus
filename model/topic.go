package model

import "time"

type Topic struct {
	ID        int64
	BoardID   int64
	IsSticky  bool
	IsLocked  bool
	Posts     int64
	UpdatedAt time.Time
	Post      Post
}

// TopicRequest represents the request to create or update a topic.
type TopicRequest struct {
	BoardID  int64
	IsSticky bool
	IsLocked bool
	Subject  string
	Content  string
}
