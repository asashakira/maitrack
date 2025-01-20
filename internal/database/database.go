package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Connect(port, dbURL string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	maxRetries := 5
	retryDelay := 2 * time.Second

	for attempts := 1; attempts <= maxRetries; attempts++ {
		pool, err = pgxpool.New(context.Background(), dbURL)
		if err == nil {
			// Successfully connected
			fmt.Println("Successfully connected to the database!")
			break
		}

		// If connection fails, log the error and retry after a delay
		log.Printf("Attempt %d/%d: Can't connect to database, retrying in %s...\n", attempts, maxRetries, retryDelay)
		time.Sleep(retryDelay)
	}

	// If the connection is still unsuccessful after retries, exit
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	if err := runMigrations(pool); err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	return pool, nil
}

// goose migrations
func runMigrations(pool *pgxpool.Pool) error {
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "./internal/database/migrations"); err != nil {
		return err
	}
	return nil
}
