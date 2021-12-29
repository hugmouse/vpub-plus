package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

const schemaVersion = 1

func Migrate(db *sql.DB) {
	var currentVersion int
	db.QueryRow(`SELECT version FROM schema_version`).Scan(&currentVersion)

	fmt.Println("Current schema version:", currentVersion)
	fmt.Println("Latest schema version:", schemaVersion)

	for version := currentVersion + 1; version <= schemaVersion; version++ {
		fmt.Println("Migrating to version:", version)

		tx, err := db.Begin()
		if err != nil {
			log.Fatal("[Migrate] ", err)
		}

		rawSQL := SqlMap["schema_version_"+strconv.Itoa(version)]
		if rawSQL == "" {
			log.Fatalf("[Migrate] missing migration %d", version)
		}
		_, err = tx.Exec(string(rawSQL))
		if err != nil {
			tx.Rollback()
			log.Fatal("[Migrate] ", err)
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			tx.Rollback()
			log.Fatal("[Migrate] ", err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_version (version) VALUES ($1)`, version); err != nil {
			tx.Rollback()
			log.Fatal("[Migrate] ", err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatal("[Migrate] ", err)
		}
	}
}
