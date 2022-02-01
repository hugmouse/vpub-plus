package model

type Forum struct {
	Id       int64
	Name     string
	Boards   []Board
	Position int64
	IsLocked bool
}

// ForumRequest represents the request to create or update a forum.
type ForumRequest struct {
	Name     string
	Position int64
	IsLocked bool
}
