package model

import "time"

type Board struct {
	Id          int64
	Name        string
	Description string
	Topics      int64
	Posts       int64
	UpdatedAt   time.Time
}
