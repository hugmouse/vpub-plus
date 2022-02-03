package model

import "time"

type Key struct {
	Id        int64
	Key       string
	CreatedAt time.Time
}
