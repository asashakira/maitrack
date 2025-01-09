package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/asashakira/mai.gg-api/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	NoteDesigner  pgtype.Text      `json:"noteDesigner"`
	MaxDxScore    int32            `json:"maxDxScore"`
	IsValid       bool             `json:"isValid"`
	UpdatedAt     pgtype.Timestamp `json:"updatedAt"`
	CreatedAt     pgtype.Timestamp `json:"createdAt"`
}

func (h *Handler) GetAllBeatmaps(w http.ResponseWriter, r *http.Request) {
	beatmaps, err := h.queries.GetAllBeatmaps(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("GetAllBeatmaps %v", err))
		return
	}
	respondWithJSON(w, 200, ConvertBeatmaps(beatmaps))
}

func (h *Handler) GetBeatmapsBySongID(w http.ResponseWriter, r *http.Request) {
	songID, err := uuid.Parse(chi.URLParam(r, "songID"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing songID: %v", err))
		return
	}

	beatmaps, err := h.queries.GetBeatmapsBySongID(r.Context(), songID)
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No beatmaps found with provided songID: %s", err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetBeatmapsBySongID %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	// log.Println("GetSongBySongID:", ConvertBeatmaps(beatmaps))
	respondWithJSON(w, 200, ConvertBeatmaps(beatmaps))
}

func (h *Handler) CreateBeatmap(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		SongID        uuid.UUID      `json:"songID"`
		Difficulty    string         `json:"difficulty"`
		Level         string         `json:"level"`
		InternalLevel pgtype.Numeric `json:"internalLevel"`
		Type          string         `json:"type"`
		TotalNotes    int32          `json:"totalNotes"`
		Tap           int32          `json:"tap"`
		Hold          int32          `json:"hold"`
		Slide         int32          `json:"slide"`
		Touch         int32          `json:"touch"`
		Break         int32          `json:"break"`
		NoteDesigner  pgtype.Text    `json:"noteDesigner"`
		MaxDxScore    int32          `json:"maxDxScore"`
		PlayCount     int32          `json:"playCount"`
		IsValid       bool           `json:"isValid"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	beatmap, err := h.queries.CreateBeatmap(r.Context(), database.CreateBeatmapParams{
		BeatmapID:     uuid.New(),
		SongID:        params.SongID,
		Difficulty:    params.Difficulty,
		Level:         params.Level,
		InternalLevel: params.InternalLevel,
		Type:          params.Type,
		TotalNotes:    params.TotalNotes,
		Tap:           params.Tap,
		Hold:          params.Hold,
		Slide:         params.Slide,
		Touch:         params.Touch,
		Break:         params.Break,
		NoteDesigner:  params.NoteDesigner,
		MaxDxScore:    params.MaxDxScore,
		IsValid:       params.IsValid,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("CreateBeatmap %v", err))
		return
	}
	respondWithJSON(w, 200, ConvertBeatmap(beatmap))
}

func (h *Handler) UpdateBeatmap(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		BeatmapID     uuid.UUID       `json:"beatmapID"`
		SongID        *uuid.UUID      `json:"songID,omitempty"`
		Difficulty    string          `json:"difficulty,omitempty"`
		Level         string          `json:"level,omitempty"`
		InternalLevel *pgtype.Numeric `json:"internalLevel,omitempty"`
		Type          string          `json:"type,omitempty"`
		TotalNotes    *int32          `json:"totalNotes,omitempty"`
		Tap           *int32          `json:"tap,omitempty"`
		Hold          *int32          `json:"hold,omitempty"`
		Slide         *int32          `json:"slide,omitempty"`
		Touch         *int32          `json:"touch,omitempty"`
		Break         *int32          `json:"break,omitempty"`
		NoteDesigner  *pgtype.Text    `json:"noteDesigner,omitempty"`
		MaxDxScore    *int32          `json:"maxDxScore,omitempty"`
		PlayCount     *int32          `json:"playCount,omitempty"`
		IsValid       *bool           `json:"isValid,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// Fetch existing beatmap
	beatmap, err := h.queries.GetBeatmapByBeatmapID(r.Context(), params.BeatmapID)
	if err != nil {
		errorMessage := fmt.Sprintf("beatmap not found %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	// // Parse InternalLevel
	// var internalLevel pgtype.Numeric
	// if err := internalLevel.Scan(params.InternalLevel); err != nil {
	// 	respondWithError(w, 400, fmt.Sprintf("Invalid InternalLevel: %v", err))
	// 	return
	// }

	updatedBeatmap, err := h.queries.UpdateBeatmap(r.Context(), database.UpdateBeatmapParams{
		BeatmapID:     params.BeatmapID,
		SongID:        ifNotNil(params.SongID, beatmap.SongID),
		Difficulty:    ifNotEmpty(params.Difficulty, beatmap.Difficulty),
		Level:         ifNotEmpty(params.Level, beatmap.Level),
		InternalLevel: ifNotNil(params.InternalLevel, beatmap.InternalLevel),
		Type:          ifNotEmpty(params.Type, beatmap.Type),
		TotalNotes:    ifNotNil(params.TotalNotes, beatmap.TotalNotes),
		Tap:           ifNotNil(params.Tap, beatmap.Tap),
		Hold:          ifNotNil(params.Hold, beatmap.Hold),
		Slide:         ifNotNil(params.Slide, beatmap.Slide),
		Touch:         ifNotNil(params.Touch, beatmap.Touch),
		Break:         ifNotNil(params.Break, beatmap.Break),
		NoteDesigner:  ifNotNil(params.NoteDesigner, beatmap.NoteDesigner),
		MaxDxScore:    ifNotNil(params.MaxDxScore, beatmap.MaxDxScore),
		IsValid:       ifNotNil(params.IsValid, beatmap.IsValid),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("UpdateBeatmap %v", err))
		return
	}
	respondWithJSON(w, 200, ConvertBeatmap(updatedBeatmap))
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
		IsValid:       dbBeatmap.IsValid,
		UpdatedAt:     dbBeatmap.UpdatedAt,
		CreatedAt:     dbBeatmap.CreatedAt,
	}
}
