package main

import (
	"log"
	"os"

	"github.com/asashakira/mai.gg/internal/api"
	"github.com/asashakira/mai.gg/internal/cron"
	"github.com/asashakira/mai.gg/internal/database"
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

	// run database migration
	if err := database.Migrate(pool); err != nil {
		log.Printf("migration error: %s", err)
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
