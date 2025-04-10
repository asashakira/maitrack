package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	database "github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/pkg/maimaiclient"
	"github.com/asashakira/maitrack/internal/utils"
)

// scrape rating and playcounts from maimaidxnet
func ScrapeUserData(segaID, segaPassword string) (database.UserDatum, error) {
	// Login
	m := maimaiclient.New()
	err := m.Login(segaID, segaPassword)
	if err != nil {
		return database.UserDatum{}, fmt.Errorf("failed to login to maimai: %w", err)
	}

	// Fetch playerData page
	r, err := m.HTTPClient.Get(maimaiclient.BaseURL + "/playerData")
	if err != nil {
		return database.UserDatum{}, fmt.Errorf("request failed: %w", err)
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return database.UserDatum{}, err
	}

	// rating
	rating, atoiErr := utils.StringToInt32(doc.Find(".rating_block").Text())
	if atoiErr != nil {
		return database.UserDatum{}, atoiErr
	}

	// play count
	playCounts := strings.Split(doc.Find(".m_5.m_b_5.t_r.f_12").Text(), "：")
	seasonPlayCount, atoiErr := utils.StringToInt32(utils.RemoveFromString(playCounts[1], `[^\d+]`))
	if atoiErr != nil {
		return database.UserDatum{}, fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}
	TotalPlayCount, atoiErr := utils.StringToInt32(utils.RemoveFromString(playCounts[2], `[^\d+]`))
	if atoiErr != nil {
		return database.UserDatum{}, fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}

	userData := database.UserDatum{
		Rating: rating,
		SeasonPlayCount: seasonPlayCount,
		TotalPlayCount: TotalPlayCount,
	}

	return userData, nil
}
