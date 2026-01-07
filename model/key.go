package model

import "time"

type Key struct {
	ID        int64
	Key       string
	CreatedAt time.Time
}
