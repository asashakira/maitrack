package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database "github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) CreateScore(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		BeatmapID     string `json:"beatmapID"`
		SongID        string `json:"songID"`
		UserID        string `json:"userID"`
		Accuracy      string `json:"accuracy"`
		MaxCombo      int32  `json:"maxCombo"`
		DxScore       int32  `json:"dxScore"`
		TapCritical   int32  `json:"tapCritical"`
		TapPerfect    int32  `json:"tapPerfect"`
		TapGreat      int32  `json:"tapGreat"`
		TapGood       int32  `json:"tapGood"`
		TapMiss       int32  `json:"tapMiss"`
		HoldCritical  int32  `json:"holdCritical"`
		HoldPerfect   int32  `json:"holdPerfect"`
		HoldGreat     int32  `json:"holdGreat"`
		HoldGood      int32  `json:"holdGood"`
		HoldMiss      int32  `json:"holdMiss"`
		SlideCritical int32  `json:"slideCritical"`
		SlidePerfect  int32  `json:"slidePerfect"`
		SlideGreat    int32  `json:"slideGreat"`
		SlideGood     int32  `json:"slideGood"`
		SlideMiss     int32  `json:"slideMiss"`
		TouchCritical int32  `json:"touchCritical"`
		TouchPerfect  int32  `json:"touchPerfect"`
		TouchGreat    int32  `json:"touchGreat"`
		TouchGood     int32  `json:"touchGood"`
		TouchMiss     int32  `json:"touchMiss"`
		BreakCritical int32  `json:"breakCritical"`
		BreakPerfect  int32  `json:"breakPerfect"`
		BreakGreat    int32  `json:"breakGreat"`
		BreakGood     int32  `json:"breakGood"`
		BreakMiss     int32  `json:"breakMiss"`
		Fast          int32  `json:"fast"`
		Late          int32  `json:"late"`
		PlayedAt      string `json:"playedAt"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	if params.BeatmapID == "" || params.SongID == "" || params.UserID == "" {
		utils.RespondWithError(w, 400, "BeatmapID, SongID, and UserID are required")
		return
	}

	playedAt, err := utils.StringToUTCTime(params.PlayedAt)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("error parsing played at date: %v", err))
		return
	}

	score, err := h.queries.CreateScore(r.Context(), database.CreateScoreParams{
		ScoreID:       uuid.New(),
		BeatmapID:     uuid.MustParse(params.BeatmapID),
		SongID:        uuid.MustParse(params.SongID),
		UserID:        uuid.MustParse(params.UserID),
		Accuracy:      params.Accuracy,
		MaxCombo:      params.MaxCombo,
		DxScore:       params.DxScore,
		TapCritical:   params.TapCritical,
		TapPerfect:    params.TapPerfect,
		TapGreat:      params.TapGreat,
		TapGood:       params.TapGood,
		TapMiss:       params.TapMiss,
		HoldCritical:  params.HoldCritical,
		HoldPerfect:   params.HoldPerfect,
		HoldGreat:     params.HoldGreat,
		HoldGood:      params.HoldGood,
		HoldMiss:      params.HoldMiss,
		SlideCritical: params.SlideCritical,
		SlidePerfect:  params.SlidePerfect,
		SlideGreat:    params.SlideGreat,
		SlideGood:     params.SlideGood,
		SlideMiss:     params.SlideMiss,
		TouchCritical: params.TouchCritical,
		TouchPerfect:  params.TouchPerfect,
		TouchGreat:    params.TouchGreat,
		TouchGood:     params.TouchGood,
		TouchMiss:     params.TouchMiss,
		BreakCritical: params.BreakCritical,
		BreakPerfect:  params.BreakPerfect,
		BreakGreat:    params.BreakGreat,
		BreakGood:     params.BreakGood,
		BreakMiss:     params.BreakMiss,
		Fast:          params.Fast,
		Late:          params.Late,
		PlayedAt:      pgtype.Timestamp{Time: playedAt, Valid: true},
	})
	if err != nil {
		errorMessage := fmt.Sprintf("CreateScore %v", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}
	utils.RespondWithJSON(w, 200, score)
}

type ScoresResponse struct {
	Scores     []database.GetScoreByMaiIDRow `json:"scores"`
	NextOffset int                           `json:"nextOffset,omitEmpty"`
	HasMore    bool                          `json:"hasMore"`
}

// gets score by maiID (gameName + tagLine)
func (h *Handler) GetScoresByMaiID(w http.ResponseWriter, r *http.Request) {
	maiID := chi.URLParam(r, "maiID")

	// get gamename and tagline
	gameName, tagLine, decodeMaiIDErr := decodeMaiID(maiID)
	if decodeMaiIDErr != nil {
		errorMessage := fmt.Sprintf("error decoding maiID %s", decodeMaiIDErr)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	limit, limitErr := strconv.Atoi(r.URL.Query().Get("limit"))
	if limitErr != nil {
		limit = 10
	}
	offset, offsetErr := strconv.Atoi(r.URL.Query().Get("offset"))
	if offsetErr != nil {
		offset = 0
	}

	scores, err := h.queries.GetScoreByMaiID(r.Context(), database.GetScoreByMaiIDParams{
		GameName: gameName,
		TagLine:  tagLine,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No score found with provided MaiID '%s': %s", maiID, err)
			log.Println(errorMessage)
			utils.RespondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetScoresByMaiID %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	nextOffset := offset + limit
	hasMore := len(scores) == limit

	response := ScoresResponse{
		Scores:     scores,
		NextOffset: nextOffset,
		HasMore:    hasMore,
	}

	utils.RespondWithJSON(w, 200, response)
}
