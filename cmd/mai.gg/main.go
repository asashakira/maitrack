package main

import (
	"log"
	"os"

	"github.com/asashakira/mai.gg-api/internal/api"
	"github.com/asashakira/mai.gg-api/internal/api/handler"
	"github.com/asashakira/mai.gg-api/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// connect to DB
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	conn, err := database.Connect(port, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(conn)
	server := api.New(h)
	err = server.Run(port)
	if err != nil {
		log.Fatalf("failed to run server: %s", err)
	}
}
