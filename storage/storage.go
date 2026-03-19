package storage

import (
	"database/sql"
	"sync"
	"time"

	"vpub/model"
)

type Storage struct {
	db               *sql.DB
	settingsCache    *model.Settings
	settingsCacheTTL time.Time
	settingsMu       sync.RWMutex
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}
