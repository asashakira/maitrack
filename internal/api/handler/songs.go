package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/asashakira/mai.gg/internal/api/model"
	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AltKey      string `json:"altkey"`
		Title       string `json:"title"`
		Artist      string `json:"artist"`
		Genre       string `json:"genre"`
		Bpm         string `json:"bpm"`
		ImageUrl    string `json:"imageUrl"`
		Version     string `json:"version"`
		IsUtage     bool   `json:"isUtage"`
		IsAvailable bool   `json:"isAvailable"`
		ReleaseDate string `json:"releaseDate"`
		DeleteDate  string `json:"deleteDate,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	if params.AltKey == "" {
		respondWithError(w, 400, fmt.Sprintf("altkey not provided %v", err))
		return
	}

	releaseDate, err := time.Parse("2006-01-02", params.ReleaseDate)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing release date: %v", err))
		return
	}

	var deleteDate pgtype.Date
	if params.DeleteDate != "" {
		parsedDate, err := time.Parse("2006-01-02", params.DeleteDate)
		deleteDate.Scan(parsedDate)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error parsing delete date: %v", err))
			return
		}
	}

	song, err := h.queries.CreateSong(r.Context(), database.CreateSongParams{
		SongID:      uuid.New(),
		AltKey:      params.AltKey,
		Title:       params.Title,
		Artist:      params.Artist,
		Genre:       params.Genre,
		Bpm:         params.Bpm,
		ImageUrl:    params.ImageUrl,
		Version:     params.Version,
		IsUtage:     params.IsUtage,
		IsAvailable: params.IsAvailable,
		ReleaseDate: pgtype.Date{Time: releaseDate, Valid: true},
		DeleteDate:  deleteDate,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("CreateSong %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, model.ConvertSong(song))
}

func (h *Handler) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := h.queries.GetAllSongs(r.Context())
	if err != nil {
		errorMessage := fmt.Sprintf("GetAllSongs %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, model.ConvertSongs(songs))
}

func (h *Handler) GetSongBySongID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing songID: %v", err))
		return
	}

	song, err := h.queries.GetSongBySongID(r.Context(), id)
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No song found with provided songID: %s", err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetSongBySongID %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, model.ConvertSong(song))
}

// get song using altkey
// altkey is made by
// combining title and artist
// all lowercase
// remove except these `[一-龠ぁ-ゔァ-ヴーa-zA-Z0-9ａ-ｚＡ-Ｚ０-９々〆〤ヶ]+`
func (h *Handler) GetSongByAltKey(w http.ResponseWriter, r *http.Request) {
	altkey := chi.URLParam(r, "altkey")
	altkey, err := url.QueryUnescape(altkey)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error decoding altkey from url: %v", err))
		return
	}

	song, err := h.queries.GetSongByAltKey(r.Context(), altkey)
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No song found with provided altkey '%s': %s", altkey, err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetSongByAltKey %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, model.ConvertSong(song))
}

// may return multiple songs
func (h *Handler) GetSongsByTitle(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	title, err := url.QueryUnescape(title)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error decoding title from url: %v", err))
		return
	}

	songs, err := h.queries.GetSongsByTitle(r.Context(), title)
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No song found with provided title '%s': %s", title, err)
			log.Println(errorMessage)
			respondWithError(w, 404, errorMessage)
			return
		}
		// Handle other errors
		errorMessage := fmt.Sprintf("GetSongByTitle %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, model.ConvertSongs(songs))
}

func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		SongID      uuid.UUID `json:"songID"`
		AltKey      *string    `json:"altkey,omitempty"`
		Title       *string    `json:"title,omitempty"`
		Artist      *string    `json:"artist,omitempty"`
		Genre       *string    `json:"genre,omitempty"`
		Bpm         *string    `json:"bpm,omitempty"`
		ImageUrl    *string    `json:"imageUrl,omitempty"`
		Version     *string    `json:"version,omitempty"`
		IsUtage     *bool     `json:"isUtage,omitempty"`
		IsAvailable *bool     `json:"isAvailable,omitempty"`
		ReleaseDate string    `json:"releaseDate,omitempty"`
		DeleteDate  string    `json:"deleteDate,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	// Fetch existing song
	song, err := h.queries.GetSongBySongID(r.Context(), params.SongID)
	if err != nil {
		errorMessage := fmt.Sprintf("song not found %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	// update release date if provided
	releaseDate := song.ReleaseDate
	if params.ReleaseDate != "" {
		parsedDate, err := time.Parse("2006-01-02", params.ReleaseDate)
		releaseDate.Scan(parsedDate)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error parsing release date: %v", err))
			return
		}
	}

	// update delete date if provided
	deleteDate := song.DeleteDate
	if params.DeleteDate != "" {
		parsedDate, err := time.Parse("2006-01-02", params.DeleteDate)
		deleteDate.Scan(parsedDate)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("error parsing delete date: %v", err))
			return
		}
	}

	// Update only the fields provided in the request
	updatedSong, err := h.queries.UpdateSong(r.Context(), database.UpdateSongParams{
		SongID:      song.SongID,
		AltKey:      ifNotNil(params.AltKey, song.AltKey),
		Title:       ifNotNil(params.Title, song.Title),
		Artist:      ifNotNil(params.Artist, song.Artist),
		Genre:       ifNotNil(params.Genre, song.Genre),
		Bpm:         ifNotNil(params.Bpm, song.Bpm),
		ImageUrl:    ifNotNil(params.ImageUrl, song.ImageUrl),
		Version:     ifNotNil(params.Version, song.Version),
		IsUtage:     ifNotNil(params.IsUtage, song.IsUtage),
		IsAvailable: ifNotNil(params.IsAvailable, song.IsAvailable),
		ReleaseDate: releaseDate,
		DeleteDate:  deleteDate,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("UpdateSong %v", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	respondWithJSON(w, 200, model.ConvertSong(updatedSong))
}

