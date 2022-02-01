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

// Patch updates forum fields.
func (fr ForumRequest) Patch(forum Forum) Forum {
	forum.Name = fr.Name
	forum.Position = fr.Position
	forum.IsLocked = fr.IsLocked
	return forum
}
