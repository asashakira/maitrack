package scraper

import (
	"log"

	database "github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ScrapeSongsAndBeatmaps(pool *pgxpool.Pool) {
	log.Println("ScrapeSongsAndBeatmaps START")

	queries := database.New(pool)

	maimaisongs, err := fetchMaimaiSongs()
	if err != nil {
		log.Printf("failed loading maimai songs: %s\n", err)
	}

	for _, ms := range maimaisongs {
		handleEdgeCases(&ms)

		song, err := upsertSong(queries, ms)
		if err != nil {
			log.Println(err)
		}

		err = handleBeatmaps(queries, song.ID, ms)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("ScrapeSongsAndBeatmaps DONE")
}
