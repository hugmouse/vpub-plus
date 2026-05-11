package model

type Group struct {
	ID          int64
	Name        string
	MemberCount int64
}

type GroupRequest struct {
	Name string
}
