package middleware

import (
	"github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Middleware struct {
	queries *sqlc.Queries
}

func New(pool *pgxpool.Pool) *Middleware {
	return &Middleware{
		queries: sqlc.New(pool),
	}
}
