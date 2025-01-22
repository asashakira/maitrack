package model

import (
	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Song struct {
	SongID      uuid.UUID        `json:"songID"`
	AltKey      string           `json:"altkey"`
	Title       string           `json:"title"`
	Artist      string           `json:"artist"`
	Genre       string           `json:"genre"`
	Bpm         string           `json:"bpm"`
	ImageUrl    string           `json:"imageUrl"`
	Version     string           `json:"version"`
	IsUtage     bool             `json:"isUtage"`
	IsAvailable bool             `json:"isAvailable"`
	ReleaseDate pgtype.Date      `json:"releaseDate"`
	DeleteDate  pgtype.Date      `json:"deleteDate"`
	UpdatedAt   pgtype.Timestamp `json:"updatedAt"`
	CreatedAt   pgtype.Timestamp `json:"createdAt"`
}

func ConvertSongs(dbSong []database.Song) []Song {
	songs := []Song{}
	for _, song := range dbSong {
		songs = append(songs, ConvertSong(song))
	}
	return songs
}

func ConvertSong(dbSong database.Song) Song {
	return Song{
		SongID:      dbSong.SongID,
		AltKey:      dbSong.AltKey,
		Title:       dbSong.Title,
		Artist:      dbSong.Artist,
		Genre:       dbSong.Genre,
		Bpm:         dbSong.Bpm,
		ImageUrl:    dbSong.ImageUrl,
		Version:     dbSong.Version,
		IsUtage:     dbSong.IsUtage,
		IsAvailable: dbSong.IsAvailable,
		ReleaseDate: dbSong.ReleaseDate,
		DeleteDate:  dbSong.DeleteDate,
		UpdatedAt:   dbSong.UpdatedAt,
		CreatedAt:   dbSong.CreatedAt,
	}
}
