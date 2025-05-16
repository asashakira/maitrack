package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/asashakira/maitrack/internal/database"
	"github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/scraper/gamerch"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// connect to database
	pool, err := database.Connect(port, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	g := gamerch.New(pool)
	if err != nil {
		log.Fatal(err)
	}
	// gamerchSongs, _, err := g.ScrapeAllSongs()
	song, _, err := g.ScrapeSong("https://gamerch.com/maimai/533541")
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(pool)
	updateDBSongWithGamerchSong(queries, song)
	// for _, song := range gamerchSongs {
	// 	updateDBSongWithGamerchSong(queries, song)
	// }

	gamerch.PrintSong(song)
	// gamerch.PrintBeatmaps(beatmaps)
}

func updateDBSongWithGamerchSong(queries *sqlc.Queries, song gamerch.Song) {
	id, _ := uuid.Parse(song.SongID)

	// update release date if provided
	var releaseDate pgtype.Date
	if song.ReleaseDate != "" {
		parsedDate, err := utils.StringToUTCDate(song.ReleaseDate)
		if err != nil {
			log.Printf("release date parse failed: %s\n", song.Title)
			return
		}
		releaseDate.Scan(parsedDate)
	}

	// update delete date if provided
	var deleteDate pgtype.Date
	if song.DeleteDate != "" {
		parsedDate, err := utils.StringToUTCDate(song.DeleteDate)
		if err != nil {
			log.Printf("delete date parse failed: %s\n", song.Title)
			return
		}
		deleteDate.Scan(parsedDate)
	}

	_, err := queries.UpdateSong(context.Background(), sqlc.UpdateSongParams{
		SongID:      id,
		Title:       song.Title,
		Artist:      song.Artist,
		Genre:       song.Genre,
		Bpm:         song.Bpm,
		ImageUrl:    song.ImageUrl,
		Version:     song.Version,
		IsUtage:     song.IsUtage,
		IsAvailable: song.IsAvailable,
		ReleaseDate: releaseDate,
		DeleteDate:  deleteDate,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("UpdateSong %v", err)
		log.Println(errorMessage)
		return
	}
}
