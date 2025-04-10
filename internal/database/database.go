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
	retryDelay := 5 * time.Second

	// Use a cancellable context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for attempts := 1; attempts <= maxRetries; attempts++ {
		pool, err = pgxpool.New(ctx, dbURL)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				// Successfully connected
				fmt.Println("Successfully connected to the database!")
				return pool, nil
			}

			// If ping fails, retry
			log.Printf("Database ping failed: %v, retrying in %s...\n", err, retryDelay)
		} else {
			log.Printf("Attempt %d/%d: Can't connect to database, retrying in %s...\n", attempts, maxRetries, retryDelay)
		}

		// retry after delay
		time.Sleep(retryDelay)
	}

	// If the connection is still unsuccessful after retries, exit
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// goose up
func Migrate(pool *pgxpool.Pool) error {
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "./internal/database/migration"); err != nil {
		return err
	}
	return nil
}
