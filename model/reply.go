package model

import "errors"

type Reply struct {
	Id        int64
	Author    string
	Content   string
	PostId    int64
	ParentId  *int64
	PostTitle string
	Comments  int
	Thread    []Reply
}

func (r Reply) Validate() error {
	if len(r.Content) == 0 {
		return errors.New("content is empty")
	}
	return nil
}
