package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asashakira/mai.gg/api"
	"github.com/asashakira/mai.gg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserData struct {
	ID              uuid.UUID        `json:"id"`
	UserID          uuid.UUID        `json:"userID"`
	Rating          int32            `json:"rating"`
	SeasonPlayCount int32            `json:"seasonPlayCount"`
	TotalPlayCount  int32            `json:"totalPlayCount"`
	CreatedAt       pgtype.Timestamp `json:"createdAt"`
}

func ConvertUserData(dbUserData database.UserDatum) UserData {
	return UserData{
		ID:              dbUserData.ID,
		UserID:          dbUserData.UserID,
		Rating:          dbUserData.Rating,
		SeasonPlayCount: dbUserData.SeasonPlayCount,
		TotalPlayCount:  dbUserData.TotalPlayCount,
		CreatedAt:       dbUserData.CreatedAt,
	}
}

func (h *Handler) InsertUserData(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID          uuid.UUID `json:"userID"`
		Rating          int32     `json:"rating"`
		SeasonPlayCount int32     `json:"seasonPlayCount"`
		TotalPlayCount  int32     `json:"totalPlayCount"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	err = h.queries.InsertUserData(r.Context(), database.InsertUserDataParams{
		ID:              uuid.New(),
		UserID:          params.UserID,
		Rating:          params.Rating,
		SeasonPlayCount: params.SeasonPlayCount,
		TotalPlayCount:  params.TotalPlayCount,
	})
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("Error Creating User: %v", err))
		return
	}
}

func (h *Handler) GetUserDataByMaiID(w http.ResponseWriter, r *http.Request) {
	user, err := h.queries.GetUserByMaiID(r.Context(), database.GetUserByMaiIDParams{
		GameName: chi.URLParam(r, "gameName"),
		TagLine:  chi.URLParam(r, "tagLine"),
	})
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("GetUserByMaiID Error: %v", err))
		return
	}

	userData, err := h.queries.GetUserDataByUserID(r.Context(), user.UserID)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("GetUserDataByUserID Error: %v", err))
		return
	}
	api.RespondWithJSON(w, 200, ConvertUserData(userData))
}
