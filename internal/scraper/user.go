package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asashakira/mai.gg/internal/api/model"
	"github.com/asashakira/mai.gg/pkg/maimaiclient"
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
	rating, atoiErr := strconv.Atoi(doc.Find(".rating_block").Text())
	if atoiErr != nil {
		return atoiErr
	}
	user.Rating = int32(rating)

	// play count
	playCounts := strings.Split(doc.Find(".m_5.m_b_5.t_r.f_12").Text(), "：")
	re := regexp.MustCompile(`[^\d+]`)
	seasonPlayCount, atoiErr := strconv.Atoi(re.ReplaceAllString(playCounts[1], ""))
	if atoiErr != nil {
		return fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}
	TotalPlayCount, atoiErr := strconv.Atoi(re.ReplaceAllString(playCounts[2], ""))
	if atoiErr != nil {
		return fmt.Errorf("failed to atoi seasonPlayCount: %w", atoiErr)
	}
	user.SeasonPlayCount = int32(seasonPlayCount)
	user.TotalPlayCount = int32(TotalPlayCount)

	return nil
}
