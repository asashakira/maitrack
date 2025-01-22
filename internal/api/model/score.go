package model

import (
	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/google/uuid"
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
	CreatedAt     pgtype.Timestamp `json:"createdAt,omitempty"`               // Timestamp not null
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
