package scraper

import (
	"context"
	"fmt"
	"log"

	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/asashakira/mai.gg/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ScrapeSongsAndBeatmaps(pool *pgxpool.Pool) {
	fmt.Println("ScrapeSongsAndBeatmaps Start")

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

		err = handleBeatmaps(queries, song.SongID, ms)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println("ScrapeSongsAndBeatmaps Done")
}

// scrapes user data and scores from maimaidxnet
// then update database
func ScrapeAllUsers(pool *pgxpool.Pool) {
	fmt.Println("ScrapeUsers Start")

	queries := database.New(pool)

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

		// scrape user data
		scrapedUserData, scrapeErr := ScrapeUserData(u.SegaID, decryptedSegaPassword)
		if scrapeErr != nil {
			log.Println(scrapeErr)
			return
		}

		_, createUserDataErr := queries.CreateUserData(context.Background(), database.CreateUserDataParams{
			ID:              uuid.New(),
			UserID:          u.UserID,
			GameName:        u.GameName,
			TagLine:         u.TagLine,
			Rating:          scrapedUserData.Rating,
			SeasonPlayCount: scrapedUserData.SeasonPlayCount,
			TotalPlayCount:  scrapedUserData.TotalPlayCount,
		})
		if createUserDataErr != nil {
			log.Printf("failed to create user data for user '%s#%s': %s", u.GameName, u.TagLine, createUserDataErr)
		}

		// scrape scores
		scores, scrapeErr := scrapeScores(queries, u.SegaID, decryptedSegaPassword, u.LastPlayedAt.Time)
		if scrapeErr != nil {
			log.Println(scrapeErr)
			return
		}

		var lastPlayedAt pgtype.Timestamp

		// update database
		for _, score := range scores {
			score.UserID = u.UserID
			createScoreErr := createScore(queries, score)
			if createScoreErr != nil {
				log.Println(createScoreErr)
				continue
			}

			// update beatmap if notes are not set
			updateBeatmapErr := updateBeatmap(queries, score)
			if updateBeatmapErr != nil {
				log.Println(updateBeatmapErr)
			}

			// update lastPlayedAt with the latest time
			if lastPlayedAt.Time.Before(score.PlayedAt.Time) {
				lastPlayedAt = score.PlayedAt
			}
		}

		if lastPlayedAt.Valid {
			// update LastPlayedAt
			_, updateErr := queries.UpdateUserMetadata(context.Background(), database.UpdateUserMetadataParams{
				UserID:       u.UserID,
				LastPlayedAt: lastPlayedAt,
			})
			if updateErr != nil {
				log.Printf("failed to update user metadata: %s\n", updateErr)
			}
		}
	}

	fmt.Println("ScrapeUsers Done")
}

// insert new score to database
func createScore(queries *database.Queries, score database.Score) error {
	_, createScoreErr := queries.CreateScore(context.Background(), database.CreateScoreParams{
		ScoreID:       uuid.New(),
		BeatmapID:     score.BeatmapID,
		SongID:        score.SongID,
		UserID:        score.UserID,
		Accuracy:      score.Accuracy,
		MaxCombo:      score.MaxCombo,
		DxScore:       score.DxScore,
		TapCritical:   score.TapCritical,
		TapPerfect:    score.TapPerfect,
		TapGreat:      score.TapGreat,
		TapGood:       score.TapGood,
		TapMiss:       score.TapMiss,
		HoldCritical:  score.HoldCritical,
		HoldPerfect:   score.HoldPerfect,
		HoldGreat:     score.HoldGreat,
		HoldGood:      score.HoldGood,
		HoldMiss:      score.HoldMiss,
		SlideCritical: score.SlideCritical,
		SlidePerfect:  score.SlidePerfect,
		SlideGreat:    score.SlideGreat,
		SlideGood:     score.SlideGood,
		SlideMiss:     score.SlideMiss,
		TouchCritical: score.TouchCritical,
		TouchPerfect:  score.TouchPerfect,
		TouchGreat:    score.TouchGreat,
		TouchGood:     score.TouchGood,
		TouchMiss:     score.TouchMiss,
		BreakCritical: score.BreakCritical,
		BreakPerfect:  score.BreakPerfect,
		BreakGreat:    score.BreakGreat,
		BreakGood:     score.BreakGood,
		BreakMiss:     score.BreakMiss,
		Fast:          score.Fast,
		Late:          score.Late,
		PlayedAt:      score.PlayedAt,
	})
	if createScoreErr != nil {
		return fmt.Errorf("failed to create score: %w", createScoreErr)
	}
	return nil
}

// update beatmap only if notes are not set
func updateBeatmap(queries *database.Queries, score database.Score) error {
	beatmap, getErr := queries.GetBeatmapByBeatmapID(context.Background(), score.BeatmapID)
	if getErr != nil {
		return fmt.Errorf("failed to get beatmap: %w", getErr)
	}

	// notes are set so don't bother updating
	if beatmap.TotalNotes != 0 {
		return nil
	}

	// get note count
	tap := score.TapCritical + score.TapPerfect + score.TapGreat + score.TapGood + score.TapMiss
	hold := score.HoldCritical + score.HoldPerfect + score.HoldGreat + score.HoldGood + score.HoldMiss
	slide := score.SlideCritical + score.SlidePerfect + score.SlideGreat + score.SlideGood + score.SlideMiss
	touch := score.TouchCritical + score.TouchPerfect + score.TouchGreat + score.TouchGood + score.TouchMiss
	br := score.BreakCritical + score.BreakPerfect + score.BreakGreat + score.BreakGood + score.BreakMiss
	totalNotes := tap + hold + slide + touch + br

	_, updateErr := queries.UpdateBeatmap(context.Background(), database.UpdateBeatmapParams{
		BeatmapID:     score.BeatmapID,
		SongID:        score.SongID,
		Difficulty:    beatmap.Difficulty,
		Level:         beatmap.Level,
		InternalLevel: beatmap.InternalLevel,
		Type:          beatmap.Type,
		TotalNotes:    totalNotes,
		Tap:           tap,
		Hold:          hold,
		Slide:         slide,
		Touch:         touch,
		Break:         br,
		MaxDxScore:    totalNotes * 3,
		NoteDesigner:  beatmap.NoteDesigner,
	})
	if updateErr != nil {
		return fmt.Errorf("failed to update beatmap: %w", updateErr)
	}

	return nil
}
