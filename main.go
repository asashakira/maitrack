package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/asashakira/mai.gg/api/handlers"
	"github.com/asashakira/mai.gg/internal/database"
	"github.com/asashakira/mai.gg/scraper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
		panic(err)
	}
	defer conn.Close(context.Background())

	db := database.New(conn)
	h := handlers.New(conn)

	go scraper.StartScraping(db)

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
	v1Router.Get("/healthz", handlers.HealthCheck)
	v1Router.Get("/err", handlers.ErrorCheck)

	// user routes
	v1Router.Get("/users", h.GetAllUsers)
	v1Router.Get("/users/by-id/{id}", h.GetUserByID)
	v1Router.Get("/users/by-mai-id/{gameName}/{tagLine}", h.GetUserByMaiID)
	v1Router.Post("/users", h.CreateUser)
	v1Router.Patch("/users", h.UpdateUser)

	// user data routes
	// v1Router.Get("/users/by-uuid/{uuid}/data", h.GetUserDataByUUID)
	v1Router.Get("/users/by-mai-id/{gameName}/{tagLine}/data", h.GetUserDataByMaiID)

	// scores
	// v1Router.Get("/users/by-mai-id/{gameName}/{tagLine}/scores", h.GetScoresByMaiID)
	// v1Router.Post("/scores", h.CreateRecord)
	// v1Router.Put("/records", apiHandler.UpdateRecord)
	
	// songs
	v1Router.Get("/songs", h.GetAllSongs)
	v1Router.Get("/songs/by-id/{id}", h.GetSongBySongID)
	v1Router.Get("/songs/by-altkey/{altkey}", h.GetSongByAltKey)
	v1Router.Post("/songs", h.CreateSong)
	v1Router.Patch("/songs", h.UpdateSong)

	// beatmaps
	v1Router.Get("/beatmaps", h.GetAllBeatmaps)
	v1Router.Post("/beatmaps", h.CreateBeatmap)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
