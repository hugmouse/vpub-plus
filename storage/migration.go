package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lib/pq"
)

const schemaVersion = 11

func Migrate(db *sql.DB) {
	var currentVersion int
	err := db.QueryRow(`SELECT version FROM schema_version`).Scan(&currentVersion)
	if err != nil {
		var postgresErr *pq.Error
		// 42P01 is PostgreSQL error code for "table does not exist"
		if errors.As(err, &postgresErr) && postgresErr.Code == "42P01" {
			log.Println("[Migrate] schema_version table doesn't exist, treating current version as 0")
		} else {
			log.Fatal("[Migrate] failed to select version from schema_version: ", err)
		}
	}

	log.Println("Current schema version:", currentVersion)
	log.Println("Latest schema version:", schemaVersion)

	for version := currentVersion + 1; version <= schemaVersion; version++ {
		fmt.Println("Migrating to version:", version)

		tx, err := db.Begin()
		if err != nil {
			log.Println("[Migrate] ", err)
			log.Println("[Migrate] rolling back transaction")
			err = tx.Rollback()
			if err != nil {
				log.Fatal("[Migrate] Failed to rollback after error in version", version, ":", err)
			}
			os.Exit(1)
		}

		rawSQL := SqlMap["schema_version_"+strconv.Itoa(version)]
		if rawSQL == "" {
			log.Fatalf("[Migrate] missing migration %d", version)
		}

		_, err = tx.Exec(rawSQL)
		if err != nil {
			log.Println("[Migrate] Error executing migration for version", version, ":", err)
			log.Println("[Migrate] rolling back transaction")
			err = tx.Rollback()
			if err != nil {
				log.Fatal("[Migrate] Failed to rollback after error in version", version, ":", err)
			}
			log.Println("[Migrate] rolled back successfully")
			os.Exit(1)
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			log.Println("[Migrate] Error deleting previous schema version record:", err)
			log.Println("[Migrate] rolling back transaction")
			err = tx.Rollback()
			if err != nil {
				log.Fatal("[Migrate] Failed to rollback after error in deletion:", err)
			}
			log.Println("[Migrate] rolled back successfully")
			os.Exit(1)
		}

		if _, err := tx.Exec(`INSERT INTO schema_version (version) VALUES ($1)`, version); err != nil {
			log.Println("[Migrate] Error inserting new schema version", version, ":", err)
			log.Println("[Migrate] rolling back transaction")
			err = tx.Rollback()
			if err != nil {
				log.Fatal("[Migrate] Failed to rollback after error in insertion:", err)
			}
			log.Println("[Migrate] rolled back successfully")
			os.Exit(1)
		}

		if err := tx.Commit(); err != nil {
			log.Fatal("[Migrate] Failed to commit migration for version", version, ":", err)
		}
	}
}
