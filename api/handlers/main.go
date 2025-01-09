package handlers

import (
	"github.com/asashakira/mai.gg/internal/database"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	queries *database.Queries
}

func New(conn *pgx.Conn) *Handler {
	return &Handler{
		queries: database.New(conn),
	}
}


