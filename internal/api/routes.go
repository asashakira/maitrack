package api

import (
	"github.com/asashakira/maitrack/internal/api/handler"
	"github.com/asashakira/maitrack/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetUpRoutes(r *chi.Mux, h *handler.Handler, m *middleware.Middleware) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://maitrack.asashakira.dev", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handler.HealthCheck) 
	v1Router.Get("/err", handler.ErrorCheck)

	// auth routes
	v1Router.Post("/auth/register", h.Register)
	v1Router.Post("/auth/login", h.Login)
	v1Router.Post("/auth/logout", h.Logout)
	v1Router.Get("/auth/me", m.Auth(h.GetMe))
	v1Router.Get("/users/healthz", m.Auth(h.GetUserHealthCheck))

	// user routes
	v1Router.Get("/users", h.GetAllUsers)
	v1Router.Get("/users/by-user-id/{userID}", h.GetUserByUserID)
	v1Router.Post("/users/by-user-id/{userID}/update", h.UpdateUserByUserID)

	// songs
	v1Router.Get("/songs", h.GetAllSongs)
	v1Router.Get("/songs/by-id/{id}", h.GetSongByID)
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
	v1Router.Get("/users/by-user-id/{userID}/scores", h.GetScoresByUserID)
	v1Router.Post("/scores", h.CreateScore)

	r.Mount("/v1", v1Router)
}
