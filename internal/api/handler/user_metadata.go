package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/asashakira/mai.gg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserMetadata struct {
	UserID       uuid.UUID        `json:"userID"`
	LastPlayedAt pgtype.Timestamp `json:"lastPlayedAt"`
	UpdatedAt    pgtype.Timestamp `json:"updatedAt"`
	CreatedAt    pgtype.Timestamp `json:"createdAt"`
}

func (h *Handler) GetUserMetadataByUserID(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "id")
	user, err := h.queries.GetUserMetadataByUserID(r.Context(), uuid.MustParse(userid))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("GetUserMetadataByUserID Error: %v", err))
		return
	}
	respondWithJSON(w, 200, ConvertUserMetadata(user))
}

func (h *Handler) UpdateUserMetadata(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID       string `json:"userID"`
		LastPlayedAt string `json:"lastPlayedAt"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	lastPlayedAt, err := utils.StringToUTCTime(params.LastPlayedAt)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing last played at date: %v", err))
		return
	}

	user, err := h.queries.UpdateUserMetadata(r.Context(), database.UpdateUserMetadataParams{
		UserID:       uuid.MustParse(params.UserID),
		LastPlayedAt: pgtype.Timestamp{Time: lastPlayedAt, Valid: true},
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("UpdateUserMetadata Error: %v", err))
		return
	}
	respondWithJSON(w, 200, ConvertUserMetadata(user))
}

func ConvertUserMetadata(db database.UserMetadatum) UserMetadata {
	return UserMetadata{
		UserID:       db.UserID,
		LastPlayedAt: db.LastPlayedAt,
		UpdatedAt:    db.UpdatedAt,
		CreatedAt:    db.CreatedAt,
	}
}
