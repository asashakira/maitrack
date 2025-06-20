package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	database "github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/asashakira/maitrack/pkg/maimaiclient"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// scrapes user data and scores from maimaidxnet
// then update database
func ScrapeAllUsers(pool *pgxpool.Pool) {
	log.Println("ScrapeAllUsers START")

	queries := database.New(pool)

	users, getErr := queries.GetAllUsers(context.Background())
	if getErr != nil {
		log.Println(getErr)
	}

	for _, u := range users {
		decryptedSegaID, decryptErr := utils.Decrypt(u.EncryptedSegaID)
		if decryptErr != nil {
			log.Printf("failed to decrypt SEGA ID: %s", decryptErr)
			return
		}
		decryptedSegaPassword, decryptErr := utils.Decrypt(u.EncryptedSegaPassword)
		if decryptErr != nil {
			log.Printf("failed to decrypt SEGA password: %s", decryptErr)
			return
		}
		m := maimaiclient.New()
		err := m.Login(decryptedSegaID, decryptedSegaPassword)
		if err != nil {
			log.Printf("failed to login to maimai with SEGA ID '%s': %s\n", u.EncryptedSegaID, err)
			return
		}
		scrapeUserErr := ScrapeUser(m, queries, ScrapeUserParams{
			ID:           u.ID,
			UserID:       u.UserID,
			LastPlayedAt: u.LastPlayedAt,
		})
		if scrapeUserErr != nil {
			log.Printf("ERROR failed to scrape user %s: %s", u.UserID, scrapeUserErr)
			return
		}
	}

	log.Println("ScrapeAllUsers DONE")
}

type ScrapeUserParams struct {
	ID           uuid.UUID
	UserID       string
	LastPlayedAt pgtype.Timestamp
}

// scrapes user data and scores from maimaidxnet
// then update database
func ScrapeUser(m *maimaiclient.Client, queries *database.Queries, user ScrapeUserParams) error {
	log.Println("START Scrape User:", user)

	// update LastScrapedAt
	_, updateErr := queries.UpdateLastScrapedAt(context.Background(), database.UpdateLastScrapedAtParams{
		UserID:        user.UserID,
		LastScrapedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
	})
	if updateErr != nil {
		log.Printf("failed to update LastScrapedAt: %s\n", updateErr)
	}

	// scrape user data and save to database
	scrapedUserData, scrapeErr := ScrapePlayerDataPage(m)
	if scrapeErr != nil {
		log.Println(scrapeErr)
		return scrapeErr
	}
	_, createUserDataErr := queries.CreateUserData(context.Background(), database.CreateUserDataParams{
		ID:              uuid.New(),
		UserUuid:        user.ID,
		Rating:          scrapedUserData.Rating,
		SeasonPlayCount: scrapedUserData.SeasonPlayCount,
		TotalPlayCount:  scrapedUserData.TotalPlayCount,
	})
	if createUserDataErr != nil {
		return fmt.Errorf("failed to create user data for '%s': %s", user.UserID, createUserDataErr)
	}

	updateProfileImageUrlErr := queries.UpdateProfileImageUrl(context.Background(), database.UpdateProfileImageUrlParams{
		UserUuid:        user.ID,
		ProfileImageUrl: scrapedUserData.ProfileImageUrl,
	})
	if updateProfileImageUrlErr != nil {
		return fmt.Errorf("failed to update profile image url: %s", updateProfileImageUrlErr)
	}

	// scrape scores
	scores, scrapeErr := scrapeScores(m, queries, user.LastPlayedAt.Time)
	if scrapeErr != nil {
		return scrapeErr
	}

	var lastPlayedAt pgtype.Timestamp

	// update database
	for _, score := range scores {
		score.UserUuid = user.ID

		// create new score
		createScoreErr := createScore(queries, score)
		if createScoreErr != nil {
			return createScoreErr
		}

		// update beatmap if notes are not set
		updateBeatmapErr := updateBeatmapNoteCounts(queries, score)
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
		_, updateErr := queries.UpdateLastPlayedAt(context.Background(), database.UpdateLastPlayedAtParams{
			UserID:       user.UserID,
			LastPlayedAt: lastPlayedAt,
		})
		if updateErr != nil {
			log.Printf("failed to update LastPlayedAt: %s\n", updateErr)
		}
	}

	log.Println("DONE Scrape User:", user)

	return nil
}

type PlayerData struct {
	Rating          int32
	SeasonPlayCount int32
	TotalPlayCount  int32
	ProfileImageUrl pgtype.Text
}

// scrape rating and playcounts from maimaidxnet
func ScrapePlayerDataPage(m *maimaiclient.Client) (PlayerData, error) {
	// Fetch playerData page
	r, err := m.HTTPClient.Get(maimaiclient.BaseURL + "/playerData")
	if err != nil {
		return PlayerData{}, fmt.Errorf("request failed: %w", err)
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return PlayerData{}, err
	}

	// profile image
	imageUrl := doc.Find("img.w_112.f_l").AttrOr(`src`, "Not Found")

	// rating
	rating, atoiErr := utils.ConvertStringToInt32(doc.Find(".rating_block").Text())
	if atoiErr != nil {
		return PlayerData{}, atoiErr
	}

	// play count
	playCounts := strings.Split(doc.Find(".m_5.m_b_5.t_r.f_12").Text(), "：")
	seasonPlayCount, atoiErr := utils.ConvertStringToInt32(utils.RemoveFromString(playCounts[1], `[^\d+]`))
	if atoiErr != nil {
		return PlayerData{}, fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}
	totalPlayCount, atoiErr := utils.ConvertStringToInt32(utils.RemoveFromString(playCounts[2], `[^\d+]`))
	if atoiErr != nil {
		return PlayerData{}, fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}

	playerData := PlayerData{
		Rating:          rating,
		SeasonPlayCount: seasonPlayCount,
		TotalPlayCount:  totalPlayCount,
		ProfileImageUrl: pgtype.Text{String: imageUrl, Valid: true},
	}

	return playerData, nil
}

func ScrapeProfileImageUrl(m *maimaiclient.Client) (string, error) {
	// Fetch playerData page
	r, err := m.HTTPClient.Get(maimaiclient.BaseURL + "/playerData")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return "", err
	}

	// profile image
	imageUrl := doc.Find("img.w_112.f_l").AttrOr(`src`, "Not Found")

	return imageUrl, nil
}

// insert new score to database
func createScore(queries *database.Queries, score database.Score) error {
	_, createScoreErr := queries.CreateScore(context.Background(), database.CreateScoreParams{
		ID:            uuid.New(),
		BeatmapID:     score.BeatmapID,
		SongID:        score.SongID,
		UserUuid:      score.UserUuid,
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
func updateBeatmapNoteCounts(queries *database.Queries, score database.Score) error {
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
		ID:            score.BeatmapID,
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
