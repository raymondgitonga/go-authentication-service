package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var fs embed.FS

type ConnectionPoolConfig struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxIdleTime time.Duration
	connMaxLifetime time.Duration
}

func NewConnectionPoolConfig() ConnectionPoolConfig {
	return ConnectionPoolConfig{
		maxOpenConns:    25,
		maxIdleConns:    25,
		connMaxIdleTime: 5 * time.Minute,
		connMaxLifetime: 5 * time.Minute,
	}
}

func NewClient(ctx context.Context, connectionDSN string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	config := NewConnectionPoolConfig()
	db.SetMaxOpenConns(config.maxOpenConns)
	db.SetMaxIdleConns(config.maxIdleConns)
	db.SetConnMaxIdleTime(config.connMaxIdleTime)
	db.SetConnMaxLifetime(config.connMaxLifetime)

	err = pingUntilAvailable(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB, dbName string) error {
	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("sourceInstance error: %w", err)
	}

	targetInstance, err := pgMigrate.WithInstance(db, new(pgMigrate.Config))
	if err != nil {
		return fmt.Errorf("targetInstance error: %w", err)
	}

	migrations, err := migrate.NewWithInstance("migrations", sourceInstance, dbName, targetInstance)
	if err != nil {
		return fmt.Errorf("migrate.NewWithInstance error: %w", err)
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate to latest version: %w", err)
	}
	return sourceInstance.Close()
}

func pingUntilAvailable(ctx context.Context, db *sql.DB) error {
	var err error
	for i := 0; i < 10; i++ {
		if err = db.PingContext(ctx); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("%w", err)
}
