package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

const schemaVersion = 8

func Migrate(db *sql.DB) {
	var currentVersion int
	err := db.QueryRow(`SELECT version FROM schema_version`).Scan(&currentVersion)
	if err != nil {
		log.Fatal("[Migrate] ", err)
	}

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
		_, err = tx.Exec(rawSQL)
		if err != nil {
			log.Println("[Migrate] ", err)
			err = tx.Rollback()
			if err != nil {
				log.Fatal("[Migrate] Failed to rollback: ", err)
			}
			os.Exit(1)
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			log.Println("[Migrate] ", err)
			err = tx.Rollback()
			if err != nil {
				log.Fatal("[Migrate] Failed to rollback: ", err)
			}
			os.Exit(1)
		}

		if _, err := tx.Exec(`INSERT INTO schema_version (version) VALUES ($1)`, version); err != nil {
			log.Println("[Migrate] ", err)
			err = tx.Rollback()
			if err != nil {
				log.Fatal("[Migrate] Failed to rollback: ", err)
			}
			os.Exit(1)
		}

		if err := tx.Commit(); err != nil {
			log.Fatal("[Migrate] ", err)
		}
	}
}
