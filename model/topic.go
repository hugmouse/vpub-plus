package model

type Topic struct {
	Id        int64
	BoardId   int64
	FirstPost Post
}
