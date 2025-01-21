package handler

import (
	"github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	queries *sqlc.Queries
}

func New(conn *pgxpool.Pool) *Handler {
	return &Handler{
		queries: sqlc.New(conn),
	}
}
