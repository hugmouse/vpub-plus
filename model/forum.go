package model

type Forum struct {
	Id       int64
	Name     string
	Boards   []Board
	Position int64
}
