package model

type Topic struct {
	Id          int64
	BoardId     int64
	Subject     string
	Author      string
	FirstPostId int64
}
