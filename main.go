package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asashakira/mai.gg-api/api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

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
		log.Fatalf("Failed to connect to database after %d attempts: %v\n", maxRetries, err)
	}

	// Make sure to close the connection when done
	defer conn.Close(context.Background())

	h := handler.New(conn)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handler.HealthCheck)
	v1Router.Get("/err", handler.ErrorCheck)

	// user routes
	v1Router.Get("/users", h.GetAllUsers)
	v1Router.Get("/users/by-id/{id}", h.GetUserByID)
	v1Router.Get("/users/by-mai-id/{gameName}/{tagLine}", h.GetUserByMaiID)
	v1Router.Get("/users/by-sega-id/{segaid}/{password}", h.GetUserBySegaID)
	v1Router.Post("/users", h.CreateUser)
	v1Router.Patch("/users", h.UpdateUser)

	// user data routes
	// v1Router.Get("/users/by-id/{id}/data", h.GetUserDataByID)
	v1Router.Get("/users/by-mai-id/{gameName}/{tagLine}/data", h.GetUserDataByMaiID)

	// user scrape metadata
	v1Router.Get("/users/by-id/{id}/metadata", h.GetUserScrapeMetadataByUserID)
	v1Router.Patch("/users/metadata", h.UpdateUserScrapeMetadata)

	// songs
	v1Router.Get("/songs", h.GetAllSongs)
	v1Router.Get("/songs/by-id/{id}", h.GetSongBySongID)
	v1Router.Get("/songs/by-altkey/{altkey}", h.GetSongByAltKey)
	v1Router.Get("/songs/by-title/{title}", h.GetSongsByTitle)
	v1Router.Post("/songs", h.CreateSong)
	v1Router.Patch("/songs", h.UpdateSong)

	// beatmaps
	v1Router.Get("/beatmaps", h.GetAllBeatmaps)
	// v1Router.Get("/beatmaps/by-id", h.GetAllBeatmaps)
	v1Router.Get("/beatmaps/by-song-id/{songID}", h.GetBeatmapsBySongID)
	v1Router.Post("/beatmaps", h.CreateBeatmap)
	v1Router.Patch("/beatmaps", h.UpdateBeatmap)

	// scores
	v1Router.Get("/scores", h.GetAllScores)
	v1Router.Get("/users/by-mai-id/{gameName}/{tagLine}/scores", h.GetScoresByMaiID)
	v1Router.Get("/users/by-id/{id}/scores", h.GetScoresByUserID)
	v1Router.Get("/scores/by-id/{id}", h.GetScoreByScoreID)
	v1Router.Post("/scores", h.CreateScore)
	// v1Router.Patch("/records", apiHandler.UpdateRecord)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v", port)
	log.Fatal(server.ListenAndServe())
}
