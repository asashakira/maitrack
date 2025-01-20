package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func Connect(port, dbURL string) (*pgx.Conn, error) {
	var conn *pgx.Conn
	var err error
	maxRetries := 5
	retryDelay := 2 * time.Second

	for attempts := 1; attempts <= maxRetries; attempts++ {
		conn, err = pgx.Connect(context.Background(), dbURL)
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

	return conn, nil
}
