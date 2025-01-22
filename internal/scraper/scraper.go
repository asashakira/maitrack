package scraper

import (
	"fmt"
	"log"

	"github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ScrapeSongsAndBeatmaps(pool *pgxpool.Pool) {
	fmt.Println("ScrapeSongsAndBeatmaps Start")

	queries := sqlc.New(pool)

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

		err = handleBeatmaps(queries, song.SongID, ms)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println("ScrapeSongsAndBeatmaps Done")
}
