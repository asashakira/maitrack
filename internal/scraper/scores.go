package scraper

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/asashakira/mai.gg/pkg/maimaiclient"
	"github.com/asashakira/mai.gg/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/schollz/progressbar/v3"
)

// scrape user scores from maimaidxnet
// returns nil if nothing to update
func scrapeScores(queries *database.Queries, segaID, segaPassword string, lastPlayedAt time.Time) ([]database.Score, error) {
	// Login
	m := maimaiclient.New()
	err := m.Login(segaID, segaPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to login to maimai: %w", err)
	}

	// Fetch records page
	r, err := m.HTTPClient.Get(maimaiclient.BaseURL + "/record")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, err
	}

	// extract hidden values from recordIDs
	var recordIDs []string
	doc.Find(`.p_10.t_l.f_0.v_b`).Each(func(i int, s *goquery.Selection) {
		// extract play time
		// skip if playedAt time is before lasyPlayedAt time
		dateStr := s.Find(`.v_b`).Text()
		playedAtString := utils.RemoveFromString(dateStr, `TRACK 0[0-9]`)
		playedAt, _ := utils.StringToUTCTime(utils.FormatDate(playedAtString))
		if !playedAt.After(lastPlayedAt) {
			return
		}

		// get hidden value for record details link
		recordID := s.Find(`input[type="hidden"]`).AttrOr("value", "")
		recordIDs = append(recordIDs, recordID)
	})

	if len(recordIDs) < 1 {
		return nil, nil
	}

	// init progressbar
	bar := progressbar.Default(int64(len(recordIDs)))

	// reverse order to insert from older scores
	slices.Reverse(recordIDs)

	// scrape scores
	var scores []database.Score
	for _, recordID := range recordIDs {
		score, err := scrapeScore(queries, m, recordID)
		if err != nil {
			log.Printf("failed scraping score: '%s' %s\n", recordID, err)
		}

		scores = append(scores, score)

		// +1 progress
		bar.Add(1)

		// wait to not get ip blocked
		time.Sleep(1 * time.Second)
	}

	return scores, nil
}

// scrape score details
// provide hiddenValue found in a hidden tag at records page
func scrapeScore(queries *database.Queries, m *maimaiclient.Client, recordID string) (database.Score, error) {
	url := maimaiclient.BaseURL + "/record/playlogDetail/?idx=" + url.QueryEscape(recordID)
	r, err := m.HTTPClient.Get(url)
	if err != nil {
		return database.Score{}, fmt.Errorf("request failed: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(r.Body)
	defer r.Body.Close()
	if err != nil {
		return database.Score{}, err
	}

	// score to return
	var score database.Score

	// accuracy
	score.Accuracy = doc.Find(`.playlog_achievement_txt`).Text()

	var comboString string
	// var syncString string
	doc.Find(`.playlog_score_block.p_5`).Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			comboString = utils.RemoveFromString(s.Text(), `[^/\d]`)
		case 1:
			// syncString = utils.RemoveFromString(s.Text(), `[^/\d]`)
		}
	})
	maxComboString := strings.Split(comboString, "/")[0]
	score.MaxCombo, _ = utils.StringToInt32(maxComboString)

	// delux score is written as "DxScore / MaxDxScore"
	dxScores := strings.Split(doc.Find(`.white.p_r_5.f_15.f_r`).Text(), "/")
	score.DxScore, _ = utils.StringToInt32(utils.RemoveFromString(dxScores[0], `[^\d]`)) // remove non numbers then convert

	// beatmap type
	// determine by if there are touch notes or not
	beatmapType := "dx"

	// note details
	doc.Find(`.playlog_notes_detail td`).Each(func(i int, s *goquery.Selection) {
		// Determine the note type and index
		noteType := ""
		var idx int
		switch i / 5 {
		case 1:
			noteType = "Tap"
			idx = i % 5
		case 2:
			noteType = "Hold"
			idx = i % 5
		case 3:
			noteType = "Slide"
			idx = i % 5
		case 4:
			if s.Text() == "" {
				// no touch notes = std map
				beatmapType = "std"
				return
			}
			noteType = "Touch"
			idx = i % 5
		case 5:
			noteType = "Break"
			idx = i % 5
		}

		setNoteValue(&score, noteType, idx, s.Text())
	})

	// fast / late
	doc.Find(`.playlog_fl_block.m_5.f_r.f_12 .p_t_5`).Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			score.Fast, _ = utils.StringToInt32(s.Text())
		case 1:
			score.Late, _ = utils.StringToInt32(s.Text())
		}
	})

	// played at
	dateStr := doc.Find(`.sub_title.t_c.f_r.f_11 .v_b`).Text()
	playedAtString := utils.RemoveFromString(dateStr, `TRACK 0[0-9]`)
	playedAt, timeParseErr := utils.StringToUTCTime(utils.FormatDate(playedAtString))
	if timeParseErr != nil {
		return database.Score{}, fmt.Errorf("failed to parse time: %w", timeParseErr)
	}
	score.PlayedAt = pgtype.Timestamp{Time: playedAt, Valid: true}

	// Title
	title := doc.Find(`.basic_block.m_5.p_5.p_l_10.f_13.break`).Text()

	// difficulty
	difficultyImgSrc := doc.Find(`.playlog_top_container.p_r img.playlog_diff.v_b`).AttrOr(`src`, "Not Found")
	difficulty := getDifficultyFromImgSrc(difficultyImgSrc)

	// imageURL
	imageURL := doc.Find(`img.music_img`).AttrOr(`src`, "Not Found")

	// ids
	songID, beatmapID, err := getSongAndBeatmapID(queries, title, difficulty, beatmapType, imageURL)
	if err != nil {
		return database.Score{}, err
	}
	score.SongID = songID
	score.BeatmapID = beatmapID

	return score, nil
}

// Helper function to set note values based on index
func setNoteValue(score *database.Score, noteType string, idx int, value string) {
	// A map for each note type with corresponding fields
	noteFieldMap := map[string][]*int32{
		"Tap":   {&score.TapCritical, &score.TapPerfect, &score.TapGreat, &score.TapGood, &score.TapMiss},
		"Hold":  {&score.HoldCritical, &score.HoldPerfect, &score.HoldGreat, &score.HoldGood, &score.HoldMiss},
		"Slide": {&score.SlideCritical, &score.SlidePerfect, &score.SlideGreat, &score.SlideGood, &score.SlideMiss},
		"Touch": {&score.TouchCritical, &score.TouchPerfect, &score.TouchGreat, &score.TouchGood, &score.TouchMiss},
		"Break": {&score.BreakCritical, &score.BreakPerfect, &score.BreakGreat, &score.BreakGood, &score.BreakMiss},
	}

	// Check if the noteType exists in the map
	if fields, ok := noteFieldMap[noteType]; ok && idx < len(fields) {
		// Convert value to integer and assign to the corresponding field
		if num, err := strconv.Atoi(value); err == nil {
			*fields[idx] = int32(num)
		}
	}
}

// get song and beatmap from database to get their ids
func getSongAndBeatmapID(queries *database.Queries, title, difficulty, beatmapType string, imageURL string) (uuid.UUID, uuid.UUID, error) {
	var songID uuid.UUID
	var beatmapID uuid.UUID

	// find song and beatmap that matches
	songs, getSongErr := queries.GetSongsByTitle(context.Background(), title)
	if getSongErr != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("song not found: %w", getSongErr)
	}

	// only song title = 'Link' returns two songs
	// for now...
	for _, s := range songs {
		// determine by imageURL
		if s.ImageUrl != imageURL {
			continue
		}

		// get beatmap
		beatmap, getBeatmapErr := queries.GetBeatmapBySongIDDifficultyAndType(context.Background(), database.GetBeatmapBySongIDDifficultyAndTypeParams{
			SongID:     s.SongID,
			Difficulty: difficulty,
			Type:       beatmapType,
		})
		if getBeatmapErr != nil {
			return uuid.Nil, uuid.Nil, fmt.Errorf("beatmap not found: %w", getBeatmapErr)
		}
		songID = beatmap.SongID
		beatmapID = beatmap.BeatmapID
	}
	return songID, beatmapID, nil
}

// parse imgSrc to determine difficulty
func getDifficultyFromImgSrc(imgSrc string) string {
	imgName := utils.RemoveFromString(imgSrc, `https://maimaidx.jp/maimai-mobile/img/`)
	switch imgName {
	case "diff_basic.png":
		return "basic"
	case "diff_advanced.png":
		return "advanced"
	case "diff_expert.png":
		return "expert"
	case "diff_master.png":
		return "master"
	case "diff_remaster.png":
		return "remaster"
	default:
		return ""
	}
}

// if the beatmap has Touch notes -> dx beatmap
func determineBeatmapType(headerText string) string {
	if strings.Contains(headerText, "Touch") {
		return "dx"
	}
	return "std"
}
