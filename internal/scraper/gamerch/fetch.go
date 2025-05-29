package gamerch

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func FetchDocumentWithRetry(url string) (*goquery.Document, error) {
	for {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0")
		req.Header.Set("Referer", "https://gamerch.com/maimai/545589")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			return goquery.NewDocumentFromReader(res.Body)
		default:
			// bodyBytes, err := io.ReadAll(res.Body)
			// if err != nil {
			// 	return nil, fmt.Errorf("error reading the body: %w", err)
			// }
			retryAfter := 30 * time.Second
			fmt.Printf("error code %d: retrying after %s seconds\n", res.StatusCode, retryAfter)

			// wait to not get ip blocked
			time.Sleep(retryAfter)
		}
	}
}

func FetchURLsFromGamerch() ([]string, error) {
	songListURL := "https://gamerch.com/maimai/545589"
	doc, err := FetchDocumentWithRetry(songListURL)
	if err != nil {
		return nil, err
	}

	var songURLs []string
	doc.Find(".markup.mu .mu__list--1").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Find("a").Attr("href")
		if !exists {
			fmt.Printf("Could not find url: %s\n", s.Text())
			return
		}
		songURLs = append(songURLs, url)
	})

	return songURLs, nil
}

// fetch deleted songs title and artist
func FetchDeletedSongs() ([]Song, error) {
	deletedSongsList := "https://gamerch.com/maimai/533442"
	doc, err := FetchDocumentWithRetry(deletedSongsList)
	if err != nil {
		return nil, err
	}

	var songs []Song
	doc.Find(".main td.mu__table--col2").Each(func(i int, s *goquery.Selection) {
		song := Song{
			Title: s.Text(),
		}
		songs = append(songs, song)
	})

	return songs, nil
}
