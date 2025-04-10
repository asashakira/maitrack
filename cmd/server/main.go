package main

import (
	"log"
	"os"

	"github.com/asashakira/maitrack/internal/api"
	"github.com/asashakira/maitrack/internal/cron"
	"github.com/asashakira/maitrack/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// connect to database
	pool, err := database.Connect(port, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// run database migration
	if err := database.Migrate(pool); err != nil {
		log.Fatal("db migration error: ", err)
	}

	// cron worker
	cronErr := cron.Run(pool)
	if cronErr != nil {
		log.Fatalf("cron error: %s", cronErr)
	}

	// run server
	server := api.New(pool)
	err = server.Run(port)
	if err != nil {
		log.Fatalf("failed to run server: %s", err)
	}
}
