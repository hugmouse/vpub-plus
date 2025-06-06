package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
	"vpub/config"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	var maxConnections int
	row := db.QueryRow("SHOW max_connections")
	if err := row.Scan(&maxConnections); err != nil {
		return nil, fmt.Errorf("failed to read max_connections: %w", err)
	}

	var usedConnections int
	row = db.QueryRow("SELECT COUNT(*) FROM pg_stat_activity")
	if err = row.Scan(&usedConnections); err != nil {
		return nil, fmt.Errorf("failed to read pg_stat_activity: %w", err)
	}

	// TODO: for cases when we have multiple vpub instances, but only one postgres
	// we need to somehow utilize only available connections
	var superUserReserved = 3
	poolSize := maxConnections - usedConnections - superUserReserved

	db.SetMaxOpenConns(poolSize)
	db.SetMaxIdleConns(poolSize / 2)
	db.SetConnMaxLifetime(3 * time.Minute)

	log.Println("[db] Max open connections is set to", poolSize)
	log.Println("[db] Max idle connections is set to", poolSize/2)

	Migrate(db)
	return db, err
}
