package api

import (
	"fmt"
	"net/http"

	"github.com/asashakira/mai.gg-api/internal/api/handler"
	"github.com/go-chi/chi/v5"
)

type API struct {
	Router *chi.Mux
}

func New(h *handler.Handler) *API {
	router := chi.NewRouter()
	SetUpRoutes(router, h)
	return &API{
		Router: router,
	}
}

func (a *API) Run(port string) error {
	server := &http.Server{
		Handler: a.Router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server starting on port %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
