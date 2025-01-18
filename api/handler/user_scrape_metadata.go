package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asashakira/mai.gg/api"
	"github.com/asashakira/mai.gg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserScrapeMetadata struct {
	UserID       uuid.UUID        `json:"userID"`
	LastPlayedAt pgtype.Timestamp `json:"lastPlayedAt"`
	UpdatedAt    pgtype.Timestamp `json:"updatedAt"`
	CreatedAt    pgtype.Timestamp `json:"createdAt"`
}

func (h *Handler) GetUserScrapeMetadataByUserID(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "id")
	user, err := h.queries.GetUserScrapeMetadataByUserID(r.Context(), uuid.MustParse(userid))
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("GetUserScrapeMetadataByUserID Error: %v", err))
		return
	}
	api.RespondWithJSON(w, 200, ConvertUserScrapeMetadata(user))
}

func (h *Handler) UpdateUserScrapeMetadata(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID       string `json:"userID"`
		LastPlayedAt string `json:"lastPlayedAt"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	lastPlayedAt, err := time.Parse("2006-01-02 15:04", params.LastPlayedAt)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("error parsing last played at date: %v", err))
		return
	}

	user, err := h.queries.UpdateUserScrapeMetadata(r.Context(), database.UpdateUserScrapeMetadataParams{
		UserID:       uuid.MustParse(params.UserID),
		LastPlayedAt: pgtype.Timestamp{Time: lastPlayedAt, Valid: true},
	})
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("UpdateUserScrapeMetadata Error: %v", err))
		return
	}
	api.RespondWithJSON(w, 200, ConvertUserScrapeMetadata(user))
}

func ConvertUserScrapeMetadata(db database.UserScrapeMetadatum) UserScrapeMetadata {
	return UserScrapeMetadata{
		UserID:       db.UserID,
		LastPlayedAt: db.LastPlayedAt,
		UpdatedAt:    db.UpdatedAt,
		CreatedAt:    db.CreatedAt,
	}
}
