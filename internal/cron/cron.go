package cron

import (
	"github.com/asashakira/mai.gg/internal/scraper"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

func Run(pool *pgxpool.Pool) error {
	var err error

	c := cron.New()

	// ScrapeSongsAndBeatmaps
	// Everyday At 1:00
	_, err = c.AddFunc("0 1 * * *", func() {
		scraper.ScrapeSongsAndBeatmaps(pool)
	})
	if err != nil {
		return err
	}

	// Scrape user data and scores
	// Everyday At 3:00
	_, err = c.AddFunc("0 3 * * *", func() {
		scraper.ScrapeAllUsers(pool)
	})
	if err != nil {
		return err
	}

	c.Start()

	// run once immediately
	scraper.ScrapeSongsAndBeatmaps(pool)
	go scraper.ScrapeAllUsers(pool)

	return nil
}
