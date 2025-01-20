package handler

import (
	"github.com/asashakira/mai.gg-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	queries *sqlc.Queries
}

func New(conn *pgx.Conn) *Handler {
	return &Handler{
		queries: sqlc.New(conn),
	}
}
