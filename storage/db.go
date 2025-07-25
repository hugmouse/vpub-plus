package storage

import (
	"database/sql"
	"fmt"
	"log"
	"vpub/config"

	_ "github.com/lib/pq"
)

const (
	// TODO: replace with logger prefix instead, need to do that for every other occurrence
	logPrefix              = "[db] "
	logInitDbConn          = logPrefix + "Initializing database connection..."
	logPingDb              = logPrefix + "Pinging database..."
	logPingSuccess         = logPrefix + "Database ping successful."
	logRunMigrations       = logPrefix + "Running database migrations..."
	logMigrationsComplete  = logPrefix + "Database migrations complete."
	logAutoConfigStart     = logPrefix + "Dynamically determining max_connections..."
	logCalcUsedConns       = logPrefix + "Calculating currently used connections from pg_stat_activity..."
	logAutoConfigWarn      = logPrefix + "Warning: database connections will be now automatically configured. To configure it yourself set the following environment variables: POSTGRES_MAX_OPEN_CONNECTIONS, POSTGRES_MAX_IDLE_CONNECTIONS, POSTGRES_MAX_LIFE_TIME"
	logErrorOpeningDb      = logPrefix + "Error opening database: %v"
	logErrorPingDb         = logPrefix + "Failed to ping DB: %v"
	logUsingMaxOpenConns   = logPrefix + "Using env POSTGRES_MAX_OPEN_CONNECTIONS: %d"
	logUsingMaxIdleConns   = logPrefix + "Using env POSTGRES_MAX_IDLE_CONNECTIONS: %d"
	logUsingMaxLifetime    = logPrefix + "Using env POSTGRES_MAX_LIFE_TIME: %s"
	logErrorAutoConfig     = logPrefix + "Errored during automatic database configuration: %v"
	logErrorReadMaxConns   = logPrefix + "Failed to read max_connections: %v"
	logServerMaxConns      = logPrefix + "PostgreSQL server max_connections: %d"
	logErrorReadPgStat     = logPrefix + "Failed to read pg_stat_activity: %v"
	logUsedPgConns         = logPrefix + "Currently used PostgreSQL connections: %d"
	logSuperuserReserved   = logPrefix + "Superuser reserved connections: %d"
	logPoolSizeNonPositive = logPrefix + "Calculated pool size is non-positive (%d). Setting to 2 to prevent errors."
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	log.Println(logInitDbConn)
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Printf(logErrorOpeningDb, err)
		return nil, err
	}

	log.Println(logPingDb)
	if err := db.Ping(); err != nil {
		log.Printf(logErrorPingDb, err)
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}
	log.Println(logPingSuccess)

	if cfg.DatabaseMaxOpenConnections > 0 && cfg.DatabaseMaxIdleConnections > 0 && cfg.DatabaseMaxLifetime > 0 {
		log.Printf(logUsingMaxOpenConns, cfg.DatabaseMaxOpenConnections)
		log.Printf(logUsingMaxIdleConns, cfg.DatabaseMaxIdleConnections)
		log.Printf(logUsingMaxLifetime, cfg.DatabaseMaxLifetime)
		db.SetMaxOpenConns(cfg.DatabaseMaxOpenConnections)
		db.SetMaxIdleConns(cfg.DatabaseMaxIdleConnections)
		db.SetConnMaxLifetime(cfg.DatabaseMaxLifetime)
	} else {
		log.Println(logAutoConfigWarn)
		poolSize, err := automaticallyConfigureConnections(db)
		if err != nil {
			log.Fatalf(logErrorAutoConfig, err)
		}
		db.SetMaxOpenConns(poolSize)
		db.SetMaxIdleConns(poolSize)
		db.SetConnMaxLifetime(cfg.DatabaseMaxLifetime)
	}

	log.Println(logRunMigrations)
	Migrate(db)
	log.Println(logMigrationsComplete)

	return db, err
}

func automaticallyConfigureConnections(db *sql.DB) (int, error) {
	var (
		maxConnections  int
		usedConnections int
		poolSize        int
		err             error
	)

	log.Println(logAutoConfigStart)
	row := db.QueryRow("SHOW max_connections")
	if err := row.Scan(&maxConnections); err != nil {
		log.Printf(logErrorReadMaxConns, err)
		return poolSize, fmt.Errorf("failed to read max_connections: %w", err)
	}
	log.Printf(logServerMaxConns, maxConnections)

	log.Println(logCalcUsedConns)
	row = db.QueryRow("SELECT COUNT(*) FROM pg_stat_activity")
	if err = row.Scan(&usedConnections); err != nil {
		log.Printf(logErrorReadPgStat, err)
		return poolSize, fmt.Errorf("failed to read pg_stat_activity: %w", err)
	}
	log.Printf(logUsedPgConns, usedConnections)

	// TODO: for cases when we have multiple vpub instances, but only one postgres
	// we need to somehow utilize only available connections
	var superUserReserved = 3
	log.Printf(logSuperuserReserved, superUserReserved)
	poolSize = maxConnections - usedConnections - superUserReserved
	if poolSize < 1 {
		log.Printf(logPoolSizeNonPositive, poolSize)
		poolSize = 2
	}

	return poolSize, nil
}
