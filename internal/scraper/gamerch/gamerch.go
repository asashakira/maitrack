package gamerch

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Gamerch struct {
	baseURL string
	queries *sqlc.Queries
}

func New(pool *pgxpool.Pool) *Gamerch {
	return &Gamerch{
		baseURL: "https://gamerch.com/maimai/",
		queries: sqlc.New(pool),
	}
}

func (g *Gamerch) ScrapeAllSongs() ([]Song, []Beatmap, error) {
	// get song urls from gamerch
	songURLs, fetchSongErr := FetchURLsFromGamerch()
	if fetchSongErr != nil {
		return []Song{}, []Beatmap{}, fmt.Errorf("%w", fetchSongErr)
	}

	// actually scrape
	// bar := progressbar.Default(int64(len(songURLs)))
	songs := []Song{}
	beatmaps := []Beatmap{}
	for _, url := range songURLs {
		s, beatmapSet, err := g.ScrapeSong(url)
		if err != nil {
			return []Song{}, []Beatmap{}, fmt.Errorf("%w", err)
		}
		songs = append(songs, s)
		beatmaps = append(beatmaps, beatmapSet...)

		// +1 progress
		// bar.Add(1)
	}
	return songs, beatmaps, nil
}

// ScrapeSong scrapes song data from a given gamerch page
// url: gamerch song url (ex: oshama -> https://gamerch.com/maimai/533541)
func (g *Gamerch) ScrapeSong(url string) (Song, []Beatmap, error) {
	// check if local file exists
	filename := utils.RemoveFromString(url, g.baseURL) + ".html"
	directory := "./tmp/html/"
	filepath := directory + filename
	exists, err := FileExists(filepath)
	if err != nil {
		return Song{}, []Beatmap{}, fmt.Errorf("%w", err)
	}
	// if not, get it
	if !exists {
		// load gamerch song page then save to ./tmp/html/
		doc, err := FetchDocumentWithRetry(url)
		if err != nil {
			return Song{}, []Beatmap{}, fmt.Errorf("%w", err)
		}
		html, _ := doc.Find(".markup.mu").Html()

		err = SaveHTMLToFile(html, directory, filename)
		if err != nil {
			return Song{}, []Beatmap{}, fmt.Errorf("%w", err)
		}

		// wait to not get ip blocked
		time.Sleep(1 * time.Second)
	}

	// load page as goquery.Document
	doc, err := LoadHTMLDocument(filepath)
	if err != nil {
		return Song{}, []Beatmap{}, fmt.Errorf("%w", err)
	}

	s, beatmapSet, err := g.parseGamerchData(doc)
	if err != nil {
		return Song{}, []Beatmap{}, fmt.Errorf("error parsing page %s: %w", url, err)
	}

	return s, beatmapSet, nil
}

// Parse the HTML Document
// Each document contains data for a single song with their beatmaps
func (g *Gamerch) parseGamerchData(doc *goquery.Document) (Song, []Beatmap, error) {
	var song Song
	var beatmaps []Beatmap

	// parse each table
	doc.Find("table").Each(func(j int, table *goquery.Selection) {
		// check top left cell to determine which table
		topLeftCell := table.Find(".mu__table--row1 .mu__table--col1").Text()

		switch {
		// 基本データ(song data table)
		case topLeftCell == "" && j == 0:
			gamerchSong, parseSongTableErr := g.parseSongTable(table)
			if parseSongTableErr != nil {
				log.Printf("parse song table error: %s", parseSongTableErr)
				return
			}

			s, err := g.getSongByDatabase(gamerchSong)
			if err != nil {
				log.Printf("getSongID error: %s", err)
				return
			}
			song = s

		// 譜面データ(beatmap data table)
		case topLeftCell == "Lv":
			b, err := g.handleBeatmapTable(table, song.SongID)
			if err != nil {
				log.Printf("handle beatmap table error: %s", err)
				return
			}
			beatmaps = append(beatmaps, b...)
		}
	})

	return song, beatmaps, nil
}

// parse gamerch song table
func (g *Gamerch) parseSongTable(table *goquery.Selection) (Song, error) {
	var genre, title, artist, bpm, releaseDate, deleteDate, version string

	table.Find(`tr`).Each(func(i int, row *goquery.Selection) {
		header := row.Find(`.mu__table--col2`).Text()
		switch header {
		case "ジャンル":
			genre = row.Find(".mu__table--col3").Text()
		case "タイトル":
			title = row.Find(".mu__table--col3").Text()
		case "アーティスト":
			artist = row.Find(".mu__table--col3").Text()
		case "BPM":
			bpm = row.Find(".mu__table--col3").Text()
		case "配信日":
			releaseDate = row.Find(".mu__table--col3").Text()
		case "削除日":
			deleteDate = row.Find(".mu__table--col3").Text()
		case "バージョン":
			version = row.Find(".mu__table--col3").Text()
		}
	})

	// ignore everything after release date
	releaseDate = FindFromString(releaseDate, `^[/0-9]+`)

	s := Song{
		AltKey:      utils.CreateAltKey(title, artist),
		Title:       title,
		Artist:      artist,
		Genre:       genre,
		Bpm:         bpm,
		Version:     version,
		ReleaseDate: releaseDate,
		DeleteDate:  deleteDate,
	}
	s.Format()
	return s, nil
}

// Gets song id from db if song exists
// if not save new song to db and return that song id
func (g *Gamerch) getSongByDatabase(gamerchSong Song) (Song, error) {
	altKey := utils.CreateAltKey(gamerchSong.Title, gamerchSong.Artist)
	s, getSongErr := g.queries.GetSongByAltKey(context.Background(), altKey)
	if getSongErr != nil {
		if strings.Contains(getSongErr.Error(), "no rows") {
			// fmt.Printf("song not found: %s %s\n", gamerchSong.Title, gamerchSong.Artist)
			return Song{}, nil
			// TODO: create new song if it does not exist in DB
			// newSong, insertErr := g.queries.CreateSong(context.Background(), sqlc.CreateSongParams{
			// 	// TODO
			// })
			// if insertErr != nil {
			// 	return "", fmt.Errorf("failed to insert song '%v': %w", gamerchSong.Title, insertErr)
			// }
			// return newSong.SongID.String(), nil
		}
		// other errors
		return Song{}, fmt.Errorf("failed to get song '%s': %w", gamerchSong.Title, getSongErr)
	}

	song := Song{
		SongID:      s.SongID.String(),
		AltKey:      s.AltKey,
		Title:       s.Title,
		Artist:      s.Artist,
		Genre:       s.Genre,
		ImageUrl:    s.ImageUrl,
		Version:     s.Version,
		IsUtage:     s.IsUtage,
		Bpm:         gamerchSong.Bpm,
		IsAvailable: gamerchSong.DeleteDate == "",
		ReleaseDate: gamerchSong.ReleaseDate,
		DeleteDate:  gamerchSong.DeleteDate,
	}
	return song, nil
}

func (g *Gamerch) handleBeatmapTable(table *goquery.Selection, songID string) ([]Beatmap, error) {
	var beatmaps []Beatmap

	// check header for missing data
	headerText := table.Find("thead th").Text()
	hasInternalLevel := strings.Contains(headerText, "定数")
	beatmapType := DetermineBeatmapType(headerText)

	// parse each row
	// each table row represents beatmap difficulty
	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		b, err := g.parseBeatmapRow(row, hasInternalLevel, beatmapType)
		if err != nil {
			// log.Printf("parse beatmap row error for song '%s': %v", s.Title, err)
		}

		if b.Difficulty == "" || b.Difficulty == "easy" || b.Difficulty == "utage" {
			// skip easy and utage maps
			// TODO: support utage maps
			return
		}

		// set SongID
		b.SongID = songID
		beatmaps = append(beatmaps, b)
	})

	return beatmaps, nil
}

// parse each row of beatmaps table
func (g *Gamerch) parseBeatmapRow(row *goquery.Selection, hasInternalLevel bool, beatmapType string) (Beatmap, error) {
	var b Beatmap
	columns := []string{"level", "internalLevel", "totalNotes", "Tap", "Hold", "Slide", "Touch", "Break"}
	current := row.Find("th")
	for _, column := range columns {
		switch column {
		case "level":
			b.Level = current.Text()
		case "internalLevel": // 譜面定数
			if !hasInternalLevel { // 定数列がなければskip
				continue
			}
			// 定数列があっても空だったら無視
			internalLevelString := current.Text()
			if internalLevelString != "" {
				internalLevel, parseErr := strconv.ParseFloat(internalLevelString, 64)
				if parseErr != nil {
					return Beatmap{}, fmt.Errorf("failed to parse internal level: %w", parseErr)
				}
				b.InternalLevel = internalLevel
			}
		case "totalNotes": // 総数
			b.TotalNotes, _ = utils.ConvertStringToInt32(current.Text())
		case "Tap":
			b.Tap, _ = utils.ConvertStringToInt32(current.Text())
		case "Hold":
			b.Hold, _ = utils.ConvertStringToInt32(current.Text())
		case "Slide":
			b.Slide, _ = utils.ConvertStringToInt32(current.Text())
		case "Touch":
			// standard譜面にはtouchない
			if beatmapType == "std" {
				continue
			}
			b.Touch, _ = utils.ConvertStringToInt32(current.Text())
		case "Break":
			b.Break, _ = utils.ConvertStringToInt32(current.Text())
		}
		current = current.Next()
	}

	b.Difficulty, _ = ParseDifficulty(row.Find("th").AttrOr("style", "yo what's up"))
	b.Type = beatmapType
	// TODO: get NoteDesigner
	b.NoteDesigner = "?"
	b.MaxDxScore = b.TotalNotes * 3

	// beatmap validation
	b.IsValid = true
	isValidBeatmap := b.TotalNotes > 0
	if !isValidBeatmap {
		b.IsValid = false
	}
	return b, nil
}
