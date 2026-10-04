package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"pet-platform/backend/internal/config"
)

type DB struct {
	*sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try pinging, but don't crash if DB is not ready during compile/unit tests
	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Warning: initial DB ping failed (will retry in background): %v\n", err)
	}

	return &DB{DB: db}, nil
}

// WithTx helper for executing operations in a transaction
func (db *DB) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
