package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asashakira/mai.gg/api"
	"github.com/asashakira/mai.gg/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Score struct {
	ScoreID   uuid.UUID        `json:"scoreID"`
	BeatmapID uuid.UUID        `json:"beatmapID"`
	SongID    uuid.UUID        `json:"songID"`
	UserID    uuid.UUID        `json:"userID"`
	Accuracy  string           `json:"accuracy"`
	MaxCombo  int32            `json:"maxCombo"`
	DxScore   int32            `json:"dxScore"`
	PlayedAt  pgtype.Timestamp `json:"playedAt"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
}

func ConvertScore(dbScore database.Score) Score {
	return Score{
		SongID:    dbScore.SongID,
		DxScore:   dbScore.DxScore,
		PlayedAt:  dbScore.PlayedAt,
		CreatedAt: dbScore.CreatedAt,
	}
}

func ConvertScores(dbScores []database.Score) []Score {
	scores := []Score{}
	for _, score := range dbScores {
		scores = append(scores, ConvertScore(score))
	}
	return scores
}

// func (h *Handler) GetScoresByMaiID(w http.ResponseWriter, r *http.Request) {
// 	gameName := chi.URLParam(r, "gameName")
// 	tagLine := chi.URLParam(r, "tagLine")
// 	records, err := h.queries.GetRecords(r.Context())
// 	if err != nil {
// 		api.RespondWithError(w, 400, fmt.Sprintf("GetScoresByMaiID Error: %v", err))
// 		return
// 	}
// 	api.RespondWithJSON(w, 200, api.ConvertRecords(records))
// }

func (h *Handler) CreateScore(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ScoreID   uuid.UUID `json:"scoreID"`
		BeatmapID uuid.UUID `json:"beatmapID"`
		SongID    uuid.UUID `json:"songID"`
		UserID    uuid.UUID `json:"userID"`
		Accuracy  string    `json:"accuracy"`
		MaxCombo  int32     `json:"maxCombo"`
		DxScore   int32     `json:"dxScore"`
		PlayedAt  string    `json:"playedAt"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	playedAtTime, err := time.Parse("2006-01-02 15:04:05.999999", params.PlayedAt)
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("Error parsing playedAt time: %v", err))
		return
	}

	_, err = h.queries.CreateScore(r.Context(), database.CreateScoreParams{
		ScoreID:   uuid.New(),
		BeatmapID: params.BeatmapID,
		SongID:    params.SongID,
		UserID:    params.UserID,
		Accuracy:  params.Accuracy,
		MaxCombo:  params.MaxCombo,
		DxScore:   params.DxScore,
		PlayedAt:  pgtype.Timestamp{Time: playedAtTime},
	})
	if err != nil {
		api.RespondWithError(w, 400, fmt.Sprintf("Error Creating Score %v", err))
		return
	}
}
