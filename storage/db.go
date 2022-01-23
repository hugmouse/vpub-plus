package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"vpub/config"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return db, err
	}
	Migrate(db)
	return db, err
}
