package handler

import (
	"fmt"
	"net/http"

	database "github.com/asashakira/mai.gg-api/internal/database/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserData struct {
	ID              uuid.UUID        `json:"id"`
	UserID          uuid.UUID        `json:"userID"`
	GameName        string           `json:"gameName"`
	TagLine         string           `json:"tagLine"`
	Rating          int32            `json:"rating"`
	SeasonPlayCount int32            `json:"seasonPlayCount"`
	TotalPlayCount  int32            `json:"totalPlayCount"`
	CreatedAt       pgtype.Timestamp `json:"createdAt"`
}

func (h *Handler) GetUserDataByMaiID(w http.ResponseWriter, r *http.Request) {
	userData, err := h.queries.GetUserDataByMaiID(r.Context(), database.GetUserDataByMaiIDParams{
		GameName: chi.URLParam(r, "gameName"),
		TagLine:  chi.URLParam(r, "tagLine"),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("GetUserDataByMaiID Error: %v", err))
		return
	}
	respondWithJSON(w, 200, ConvertUserData(userData))
}

func ConvertUserData(dbUserData database.UserDatum) UserData {
	return UserData{
		ID:              dbUserData.ID,
		UserID:          dbUserData.UserID,
		GameName:        dbUserData.GameName,
		TagLine:         dbUserData.TagLine,
		Rating:          dbUserData.Rating,
		SeasonPlayCount: dbUserData.SeasonPlayCount,
		TotalPlayCount:  dbUserData.TotalPlayCount,
		CreatedAt:       dbUserData.CreatedAt,
	}
}
