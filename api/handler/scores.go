package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asashakira/mai.gg-api/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Score struct {
	ScoreID       uuid.UUID        `json:"scoreID" db:"score_id"`             // UUID primary key
	BeatmapID     uuid.UUID        `json:"beatmapID" db:"beatmap_id"`         // UUID not null
	SongID        uuid.UUID        `json:"songID" db:"song_id"`               // UUID not null
	UserID        uuid.UUID        `json:"userID" db:"user_id"`               // UUID not null
	Accuracy      string           `json:"accuracy" db:"accuracy"`            // Text not null
	MaxCombo      int32            `json:"maxCombo" db:"max_combo"`           // Int not null
	DxScore       int32            `json:"dxScore" db:"dx_score"`             // Int not null
	TapCritical   int32            `json:"tapCritical" db:"tap_critical"`     // Int not null default 0
	TapPerfect    int32            `json:"tapPerfect" db:"tap_perfect"`       // Int not null default 0
	TapGreat      int32            `json:"tapGreat" db:"tap_great"`           // Int not null default 0
	TapGood       int32            `json:"tapGood" db:"tap_good"`             // Int not null default 0
	TapMiss       int32            `json:"tapMiss" db:"tap_miss"`             // Int not null default 0
	HoldCritical  int32            `json:"holdCritical" db:"hold_critical"`   // Int not null default 0
	HoldPerfect   int32            `json:"holdPerfect" db:"hold_perfect"`     // Int not null default 0
	HoldGreat     int32            `json:"holdGreat" db:"hold_great"`         // Int not null default 0
	HoldGood      int32            `json:"holdGood" db:"hold_good"`           // Int not null default 0
	HoldMiss      int32            `json:"holdMiss" db:"hold_miss"`           // Int not null default 0
	SlideCritical int32            `json:"slideCritical" db:"slide_critical"` // Int not null default 0
	SlidePerfect  int32            `json:"slidePerfect" db:"slide_perfect"`   // Int not null default 0
	SlideGreat    int32            `json:"slideGreat" db:"slide_great"`       // Int not null default 0
	SlideGood     int32            `json:"slideGood" db:"slide_good"`         // Int not null default 0
	SlideMiss     int32            `json:"slideMiss" db:"slide_miss"`         // Int not null default 0
	TouchCritical int32            `json:"touchCritical" db:"touch_critical"` // Int not null default 0
	TouchPerfect  int32            `json:"touchPerfect" db:"touch_perfect"`   // Int not null default 0
	TouchGreat    int32            `json:"touchGreat" db:"touch_great"`       // Int not null default 0
	TouchGood     int32            `json:"touchGood" db:"touch_good"`         // Int not null default 0
	TouchMiss     int32            `json:"touchMiss" db:"touch_miss"`         // Int not null default 0
	BreakCritical int32            `json:"breakCritical" db:"break_critical"` // Int not null default 0
	BreakPerfect  int32            `json:"breakPerfect" db:"break_perfect"`   // Int not null default 0
	BreakGreat    int32            `json:"breakGreat" db:"break_great"`       // Int not null default 0
	BreakGood     int32            `json:"breakGood" db:"break_good"`         // Int not null default 0
	BreakMiss     int32            `json:"breakMiss" db:"break_miss"`         // Int not null default 0
	Fast          int32            `json:"fast" db:"fast"`                    // Int not null default 0
	Late          int32            `json:"late" db:"late"`                    // Int not null default 0
	PlayedAt      pgtype.Timestamp `json:"playedAt" db:"played_at"`           // Timestamp not null
	CreatedAt     pgtype.Timestamp `json:"createdAt" db:"created_at"`         // Timestamp not null
}

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
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	if params.BeatmapID == "" || params.SongID == "" || params.UserID == "" {
		respondWithError(w, 400, "BeatmapID, SongID, and UserID are required")
		return
	}

	playedAt, err := time.Parse("2006-01-02 15:04", params.PlayedAt)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing played at date: %v", err))
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
		respondWithError(w, 400, errorMessage)
		return
	}

	// log.Println("CreateScore:", score)
	respondWithJSON(w, 200, score)
}

func (h *Handler) GetAllScores(w http.ResponseWriter, r *http.Request) {
	scores, err := h.queries.GetAllScores(r.Context())
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No scores found: %s", err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetAllScores %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, ConvertScores(scores))
}

func (h *Handler) GetScoreByScoreID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing songID: %v", err))
		return
	}

	score, err := h.queries.GetScoreByID(r.Context(), id)
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No score found with provided scoreID: %s", err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetScoreByID %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, ConvertScore(score))
}

func (h *Handler) GetScoresByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing userID: %s", err))
		return
	}

	scores, err := h.queries.GetScoreByUserID(r.Context(), userID)
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No score found with provided scoreID: %s", err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetScoresByUserID %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, ConvertScores(scores))
}

// gets score by maiID (gameName + tagLine)
func (h *Handler) GetScoresByMaiID(w http.ResponseWriter, r *http.Request) {
	scores, err := h.queries.GetScoreByMaiID(r.Context(), database.GetScoreByMaiIDParams{
		GameName: chi.URLParam(r, "gameName"),
		TagLine:  chi.URLParam(r, "tagLine"),
	})
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No score found with provided MaiID: %s", err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetScoresByMaiID %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, ConvertScores(scores))
}

func ConvertScore(dbScore database.Score) Score {
	return Score{
		ScoreID:       dbScore.ScoreID,
		BeatmapID:     dbScore.BeatmapID,
		SongID:        dbScore.SongID,
		UserID:        dbScore.UserID,
		Accuracy:      dbScore.Accuracy,
		MaxCombo:      dbScore.MaxCombo,
		DxScore:       dbScore.DxScore,
		TapCritical:   dbScore.TapCritical,
		TapPerfect:    dbScore.TapPerfect,
		TapGreat:      dbScore.TapGreat,
		TapGood:       dbScore.TapGood,
		TapMiss:       dbScore.TapMiss,
		HoldCritical:  dbScore.HoldCritical,
		HoldPerfect:   dbScore.HoldPerfect,
		HoldGreat:     dbScore.HoldGreat,
		HoldGood:      dbScore.HoldGood,
		HoldMiss:      dbScore.HoldMiss,
		SlideCritical: dbScore.SlideCritical,
		SlidePerfect:  dbScore.SlidePerfect,
		SlideGreat:    dbScore.SlideGreat,
		SlideGood:     dbScore.SlideGood,
		SlideMiss:     dbScore.SlideMiss,
		TouchCritical: dbScore.TouchCritical,
		TouchPerfect:  dbScore.TouchPerfect,
		TouchGreat:    dbScore.TouchGreat,
		TouchGood:     dbScore.TouchGood,
		TouchMiss:     dbScore.TouchMiss,
		BreakCritical: dbScore.BreakCritical,
		BreakPerfect:  dbScore.BreakPerfect,
		BreakGreat:    dbScore.BreakGreat,
		BreakGood:     dbScore.BreakGood,
		BreakMiss:     dbScore.BreakMiss,
		Fast:          dbScore.Fast,
		Late:          dbScore.Late,
		PlayedAt:      dbScore.PlayedAt,
		CreatedAt:     dbScore.CreatedAt,
	}
}

func ConvertScores(dbScores []database.Score) []Score {
	scores := []Score{}
	for _, score := range dbScores {
		scores = append(scores, ConvertScore(score))
	}
	return scores
}
