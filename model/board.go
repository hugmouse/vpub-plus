package model

import "time"

type Board struct {
	ID          int64
	Name        string
	Description string
	Topics      int64
	Posts       int64
	UpdatedAt   time.Time
	Position    int64
	Forum       Forum
	IsLocked    bool
}

// BoardRequest represents the request to create or update a forum.
type BoardRequest struct {
	Name        string
	Description string
	IsLocked    bool
	Position    int64
	ForumID     int64
}
