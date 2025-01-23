package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asashakira/mai.gg/internal/api/model"
	"github.com/asashakira/mai.gg/pkg/maimaiclient"
	"github.com/asashakira/mai.gg/utils"
)

// scrape rating and playcounts from maimaidxnet
func ScrapeUserData(user *model.User) error {
	// Login
	m := maimaiclient.New()
	err := m.Login(user.SegaID, user.SegaPassword)
	if err != nil {
		return fmt.Errorf("failed to login to maimai: %w", err)
	}

	// Fetch playerData page
	r, err := m.HTTPClient.Get(maimaiclient.BaseURL + "/playerData")
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return err
	}

	// rating
	rating, atoiErr := utils.StringToInt32(doc.Find(".rating_block").Text())
	if atoiErr != nil {
		return atoiErr
	}
	user.Rating = rating

	// play count
	playCounts := strings.Split(doc.Find(".m_5.m_b_5.t_r.f_12").Text(), "：")
	seasonPlayCount, atoiErr := utils.StringToInt32(utils.RemoveFromString(playCounts[1], `[^\d+]`))
	if atoiErr != nil {
		return fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}
	TotalPlayCount, atoiErr := utils.StringToInt32(utils.RemoveFromString(playCounts[2], `[^\d+]`))
	if atoiErr != nil {
		return fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}
	user.SeasonPlayCount = seasonPlayCount
	user.TotalPlayCount = TotalPlayCount

	return nil
}
