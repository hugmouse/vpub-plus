package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"vpub/config"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	log.Println("[db] Initializing database connection...")
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Printf("[db] Error opening database: %v", err)
		return nil, err
	}

	log.Println("[db] Pinging database...")
	if err := db.Ping(); err != nil {
		log.Printf("[db] Failed to ping DB: %v", err)
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}
	log.Println("[db] Database ping successful.")

	var maxConnections int
	var usedConnections int
	var poolSize int

	if cfg.DatabaseMaxOpenConnections > 0 {
		maxConnections = cfg.DatabaseMaxOpenConnections
		poolSize = maxConnections
		log.Printf("[db] Using configured POSTGRES_MAX_OPEN_CONNECTIONS: %d", maxConnections)
	} else {
		log.Println("[db] POSTGRES_MAX_OPEN_CONNECTIONS not set, dynamically determining max_connections...")
		row := db.QueryRow("SHOW max_connections")
		if err := row.Scan(&maxConnections); err != nil {
			log.Printf("[db] Failed to read max_connections: %v", err)
			return nil, fmt.Errorf("failed to read max_connections: %w", err)
		}
		log.Printf("[db] PostgreSQL server max_connections: %d", maxConnections)

		log.Println("[db] Calculating currently used connections from pg_stat_activity...")
		row = db.QueryRow("SELECT COUNT(*) FROM pg_stat_activity")
		if err = row.Scan(&usedConnections); err != nil {
			log.Printf("[db] Failed to read pg_stat_activity: %v", err)
			return nil, fmt.Errorf("failed to read pg_stat_activity: %w", err)
		}
		log.Printf("[db] Currently used PostgreSQL connections: %d", usedConnections)

		// TODO: for cases when we have multiple vpub instances, but only one postgres
		// we need to somehow utilize only available connections
		var superUserReserved = 3
		log.Printf("[db] Superuser reserved connections: %d", superUserReserved)
		poolSize = maxConnections - usedConnections - superUserReserved
		if poolSize < 1 {
			log.Printf("[db] Calculated pool size is non-positive (%d). Setting to 2 to prevent errors.", poolSize)
			poolSize = 2
		}
	}

	db.SetMaxOpenConns(poolSize)
	db.SetMaxIdleConns(poolSize / 2)
	db.SetConnMaxLifetime(cfg.DatabaseMaxLifetime)

	log.Printf("[db] Database connection settings applied: MaxOpenConns=%d, MaxIdleConns=%d, ConnMaxLifetime=%s",
		poolSize, poolSize/2, cfg.DatabaseMaxLifetime)

	log.Println("[db] Running database migrations...")
	Migrate(db)
	log.Println("[db] Database migrations complete.")

	return db, err
}
