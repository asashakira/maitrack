package api

import (
	"github.com/asashakira/mai.gg-api/internal/api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetUpRoutes(r *chi.Mux, h *handler.Handler) {
	r.Use(cors.Handler(cors.Options{
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
	v1Router.Get("/users/by-sega-id/{username}", h.GetUserByUsername)
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

	r.Mount("/v1", v1Router)
}
