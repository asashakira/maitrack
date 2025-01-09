package scraper

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Score struct {
	ScoreID   string `json:"scoreID"`
	BeatmapID string `json:"beatmapID"`
	SongID    string `json:"songID"`
	UserID    string `json:"userID"`
	Accuracy  string `json:"accuracy"`
	MaxCombo  int32  `json:"maxCombo"`
	DxScore   int32  `json:"dxScore"`
	PlayedAt  string `json:"playedAt"`
}

func scrapeScores(m *MaimaiClient) ([]Score, error) {
	res, err := m.HttpClient.Get(maimaiUrl + "/record")
	if err != nil {
		return []Score{}, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return []Score{}, err
	}

	scores := []Score{}
	doc.Find(`.p_10.t_l.f_0.v_b`).Each(func(i int, s *goquery.Selection) {
		score := Score{}
		// score.ID = s.Find(`input[type="hidden"]`).AttrOr("value", "")
		// score.Title = s.Find(`.basic_block.m_5.p_5.p_l_10.f_13.break`).Text()
		// Have to remove "TRACK 0x" from string
		dateStr := s.Find(`.v_b`).Text()
		re := regexp.MustCompile(`TRACK 0[0-9]`)
		score.PlayedAt = re.ReplaceAllString(dateStr, "")
		score.Accuracy = s.Find(`.playlog_achievement_txt`).Text()
		// delux score is written as "DxScore / MaxDxScore" so need to strip
		dxScores := strings.Split(s.Find(`.white`).Text(), "/")
		dxScore, _ := strconv.Atoi(strings.ReplaceAll(dxScores[0], " ", ""))
		// score.MaxDxScore = strings.ReplaceAll(dxScores[1], " ", "")
		// score.ImageURL = s.Find(`.music_img`).AttrOr("src", "No Image URL")
		score.DxScore = int32(dxScore)
		scores = append(scores, score)
	})

	return scores, nil
}
