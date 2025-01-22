package scraper

import (
	"context"
	"fmt"
	"log"

	"github.com/asashakira/mai.gg/internal/api/model"
	"github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/asashakira/mai.gg/utils"
	"github.com/google/uuid"
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

// scrapes user data from maimaidxnet then update database
func ScrapeAllUsers(pool *pgxpool.Pool) {
	fmt.Println("ScrapeUsers Start")

	queries := sqlc.New(pool)

	users, getErr := queries.GetAllUsers(context.Background())
	if getErr != nil {
		log.Println(getErr)
	}

	for _, u := range users {
		decryptedSegaPassword, decryptErr := utils.Decrypt(u.SegaPassword)
		if decryptErr != nil {
			log.Printf("failed to decrypt sega password: %s", decryptErr)
			break
		}

		fetchedUser := model.User{
			SegaID:       u.SegaID,
			SegaPassword: decryptedSegaPassword,
		}
		scrapeErr := ScrapeUserData(&fetchedUser)
		if scrapeErr != nil {
			log.Println(scrapeErr)
			return
		}

		_, createUserDataErr := queries.CreateUserData(context.Background(), sqlc.CreateUserDataParams{
			ID:              uuid.New(),
			UserID:          u.UserID,
			GameName:        u.GameName,
			TagLine:         u.TagLine,
			Rating:          fetchedUser.Rating,
			SeasonPlayCount: fetchedUser.SeasonPlayCount,
			TotalPlayCount:  fetchedUser.TotalPlayCount,
		})
		if createUserDataErr != nil {
			log.Println(createUserDataErr)
		}
	}

	fmt.Println("ScrapeUsers Done")
}
