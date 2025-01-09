package scraper

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/width"
)

type userData struct {
	GameName        string    `json:"gameName"`
	TagLine         string    `json:"tagLine"`
	Rating          int32     `json:"rating"`
	SeasonPlayCount int32     `json:"seasonPlayCount"`
	TotalPlayCount  int32     `json:"totalPlayCount"`
}

func scrapeUserData(m *MaimaiClient) (userData, error) {
	res, err := m.HttpClient.Get(maimaiUrl + "/playerData")
	if err != nil {
		return userData{}, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return userData{}, err
	}

	// get playcounts
	playCountRawHTML := doc.Find(".m_5").Text()
	// remove commas to make regexing easier
	playCountString := strings.ReplaceAll(playCountRawHTML, ",", "")
	re := regexp.MustCompile("[0-9]+")
	playCounts := re.FindAllString(playCountString, -1)

	rating, _ := strconv.Atoi(doc.Find(".rating_block").Text())
	seasonPlayCount, _ := strconv.Atoi(playCounts[0])
	totalPlayCount, _ := strconv.Atoi(playCounts[1])
	// TODO: add error checks for missing value
	data := userData{
		GameName:        width.Narrow.String(doc.Find(".name_block").Text()), // convert to half width
		TagLine:         "340",
		Rating:          int32(rating),
		SeasonPlayCount: int32(seasonPlayCount),
		TotalPlayCount:  int32(totalPlayCount),
	}

	return data, nil
}
