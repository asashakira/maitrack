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
	v1Router.Post("/register", m.APIKeyAuth(h.Register))

	// user routes
	v1Router.Get("/users/by-mai-id/{maiID}", m.APIKeyAuth(h.GetUserByMaiID))

	// songs
	v1Router.Get("/songs", m.APIKeyAuth(h.GetAllSongs))
	v1Router.Get("/songs/by-altkey/{altkey}", m.APIKeyAuth(h.GetSongByAltKey))
	v1Router.Get("/songs/by-title/{title}", m.APIKeyAuth(h.GetSongsByTitle))
	v1Router.Post("/songs", m.APIKeyAuth(h.CreateSong))
	v1Router.Patch("/songs", m.APIKeyAuth(h.UpdateSong))

	// beatmaps
	v1Router.Get("/beatmaps", m.APIKeyAuth(h.GetAllBeatmaps))
	v1Router.Get("/beatmaps/by-song-id/{songID}", m.APIKeyAuth(h.GetBeatmapsBySongID))
	v1Router.Post("/beatmaps", m.APIKeyAuth(h.CreateBeatmap))
	v1Router.Patch("/beatmaps", m.APIKeyAuth(h.UpdateBeatmap))

	// scores
	v1Router.Get("/users/by-mai-id/{maiID}/scores", m.APIKeyAuth(h.GetScoresByMaiID))
	v1Router.Post("/scores", m.APIKeyAuth(h.CreateScore))

	r.Mount("/v1", v1Router)
}
