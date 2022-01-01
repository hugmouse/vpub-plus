package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"vpub/config"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DatabaseFile+"?_foreign_keys=on")
	if err != nil {
		return db, err
	}
	Migrate(db)
	return db, err
}
