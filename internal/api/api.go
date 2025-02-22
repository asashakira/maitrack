package api

import (
	"log"
	"net/http"

	"github.com/asashakira/mai.gg/internal/api/handler"
	"github.com/asashakira/mai.gg/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Router *chi.Mux
}

func New(pool *pgxpool.Pool) *API {
	router := chi.NewRouter()
	h := handler.New(pool)
	m := middleware.New(pool)
	SetUpRoutes(router, h, m)
	return &API{
		Router: router,
	}
}

func (api *API) Run(port string) error {
	server := &http.Server{
		Handler: api.Router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
