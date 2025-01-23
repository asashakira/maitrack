package model

import (
	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Beatmap struct {
	BeatmapID     uuid.UUID        `json:"beatmapID"`
	SongID        uuid.UUID        `json:"songID"`
	Difficulty    string           `json:"difficulty"`
	Level         string           `json:"level"`
	InternalLevel pgtype.Numeric   `json:"internalLevel"`
	Type          string           `json:"type"`
	TotalNotes    int32            `json:"totalNotes"`
	Tap           int32            `json:"tap"`
	Hold          int32            `json:"hold"`
	Slide         int32            `json:"slide"`
	Touch         int32            `json:"touch"`
	Break         int32            `json:"break"`
	NoteDesigner  string           `json:"noteDesigner"`
	MaxDxScore    int32            `json:"maxDxScore"`
	UpdatedAt     pgtype.Timestamp `json:"updatedAt"`
	CreatedAt     pgtype.Timestamp `json:"createdAt"`
}

func ConvertBeatmaps(dbBeatmaps []database.Beatmap) []Beatmap {
	beatmaps := []Beatmap{}
	for _, beatmap := range dbBeatmaps {
		beatmaps = append(beatmaps, ConvertBeatmap(beatmap))
	}
	return beatmaps
}

func ConvertBeatmap(dbBeatmap database.Beatmap) Beatmap {
	return Beatmap{
		BeatmapID:     dbBeatmap.BeatmapID,
		SongID:        dbBeatmap.SongID,
		Difficulty:    dbBeatmap.Difficulty,
		Level:         dbBeatmap.Level,
		InternalLevel: dbBeatmap.InternalLevel,
		Type:          dbBeatmap.Type,
		TotalNotes:    dbBeatmap.TotalNotes,
		Tap:           dbBeatmap.Tap,
		Hold:          dbBeatmap.Hold,
		Slide:         dbBeatmap.Slide,
		Touch:         dbBeatmap.Touch,
		Break:         dbBeatmap.Break,
		NoteDesigner:  dbBeatmap.NoteDesigner,
		MaxDxScore:    dbBeatmap.MaxDxScore,
		UpdatedAt:     dbBeatmap.UpdatedAt,
		CreatedAt:     dbBeatmap.CreatedAt,
	}
}
