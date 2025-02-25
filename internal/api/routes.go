package api

import (
	"github.com/asashakira/mai.gg/internal/api/handler"
	"github.com/asashakira/mai.gg/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetUpRoutes(r *chi.Mux, h *handler.Handler, m *middleware.Middleware) {
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

	// auth
	v1Router.Post("/auth/register", h.Register)
    v1Router.Post("/auth/login", h.Login)
    v1Router.Get("/auth/me", m.Auth(h.GetMe))

    // user routes
    v1Router.Get("/users/by-mai-id/{maiID}", h.GetUserByMaiID)
    v1Router.Get("/users/healthz", m.Auth(h.GetUserHealthCheck))

    // songs
    v1Router.Get("/songs", h.GetAllSongs)
    v1Router.Get("/songs/by-altkey/{altkey}", h.GetSongByAltKey)
    v1Router.Get("/songs/by-title/{title}", h.GetSongsByTitle)
    v1Router.Post("/songs", h.CreateSong)
    v1Router.Patch("/songs", h.UpdateSong)

    // beatmaps
    v1Router.Get("/beatmaps", h.GetAllBeatmaps)
    v1Router.Get("/beatmaps/by-song-id/{songID}", h.GetBeatmapsBySongID)
    v1Router.Post("/beatmaps", h.CreateBeatmap)
    v1Router.Patch("/beatmaps", h.UpdateBeatmap)

    // scores
    v1Router.Get("/users/by-mai-id/{maiID}/scores", h.GetScoresByMaiID)
    v1Router.Post("/scores", h.CreateScore)

	r.Mount("/v1", v1Router)
}
