package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type maimaisong struct {
	Title      string `json:"title"`
	TitleKana  string `json:"title_kana"`
	Artist     string `json:"artist"`
	Genre      string `json:"catcode"`
	Comment    string `json:"comment"`
	Kanji      string `json:"kanji"`
	Basic      string `json:"lev_bas"`
	Advanced   string `json:"lev_adv"`
	Expert     string `json:"lev_exp"`
	Master     string `json:"lev_mas"`
	ReMaster   string `json:"lev_remas"`
	DxBasic    string `json:"dx_lev_bas"`
	DxAdvanced string `json:"dx_lev_adv"`
	DxExpert   string `json:"dx_lev_exp"`
	DxMaster   string `json:"dx_lev_mas"`
	DxReMaster string `json:"dx_lev_remas"`
	Utage      string `json:"lev_utage"`
	Version    string `json:"version"`
	Sort       string `json:"sort"`
	Date       string `json:"date"`
	Release    string `json:"release"`
	ImageUrl   string `json:"image_url"`
}

// fetch songs and beatmaps data from official api
func fetchMaimaiSongs() ([]maimaisong, error) {
	apiURL := "https://maimai.sega.jp/data/maimai_songs.json"
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0")

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error from server (status %d): %s", r.StatusCode, string(body))
	}

	var maimaisongs []maimaisong
	if err := json.Unmarshal(body, &maimaisongs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return maimaisongs, nil
}

func handleEdgeCases(ms *maimaisong) {
	// この宴譜面だけ2つあるのでコメントで差別化
	if ms.Title == "[協]青春コンプレックス" {
		suffix := "（ヒーロー級）"
		if ms.Comment == "バンドメンバーを集めて楽しもう！（入門編）" {
			suffix = "（入門編）"
		}
		ms.Title += suffix
	}

	// artist消されてる
	// 炎上したからっぽい
	if ms.Title == "ぽっぴっぽー" {
		ms.Artist = "(ラマーズP)"
	}

	// 000000 is invalid
	if ms.Release == "000000" {
		ms.Release = "060102"
	}
}

// insert if song does not exist
// update if exists
func upsertSong(queries *sqlc.Queries, ms maimaisong) (sqlc.Song, error) {
	// format releaseDate
	releaseDateString := fmt.Sprintf("20%v-%v-%v", ms.Release[0:2], ms.Release[2:4], ms.Release[4:6])
	jst, _ := time.LoadLocation("Asia/Tokyo")
	releaseDate, err := time.ParseInLocation("2006-01-02", releaseDateString, jst)
	if err != nil {
		return sqlc.Song{}, err
	}

	song, getSongErr := queries.GetSongByTitleAndArtist(context.Background(), sqlc.GetSongByTitleAndArtistParams{
		Title:  ms.Title,
		Artist: ms.Artist,
	})
	if getSongErr != nil {
		if strings.Contains(getSongErr.Error(), "no rows in result set") {
			// insert if it does not exist in DB
			newSong, createSongErr := queries.CreateSong(context.Background(), sqlc.CreateSongParams{
				SongID:      uuid.New(),
				Title:       ms.Title,
				Artist:      ms.Artist,
				Genre:       ms.Genre,
				Bpm:         "",
				ImageUrl:    ms.ImageUrl,
				Version:     versionMap[ms.Version[0:3]],
				Sort:        ms.Sort,
				IsUtage:     ms.Genre == "宴会場",
				IsAvailable: true,
				IsNew:       ms.Date == "NEW",
				ReleaseDate: pgtype.Date{Time: releaseDate, Valid: true},
			})
			if createSongErr != nil {
				return sqlc.Song{}, fmt.Errorf("failed to create song: %w", createSongErr)
			}

			uploadErr := utils.UploadImageToS3("https://maimaidx.jp/maimai-mobile/img/Music/" + ms.ImageUrl)
			if uploadErr != nil {
				return sqlc.Song{}, fmt.Errorf("failed to upload image to S3: %w", uploadErr)
			}

			// return newly created song
			return newSong, nil
		}
		// other errors
		return sqlc.Song{}, fmt.Errorf("failed to get song '%s': %w", ms.Title, getSongErr)
	}

	// update song
	updatedSong, updateErr := queries.UpdateSong(context.Background(), sqlc.UpdateSongParams{
		SongID:      song.SongID,
		Title:       ms.Title,
		Artist:      ms.Artist,
		Genre:       ms.Genre,
		Bpm:         song.Bpm,
		ImageUrl:    ms.ImageUrl,
		Version:     versionMap[ms.Version[0:3]],
		Sort:        ms.Sort,
		IsUtage:     ms.Genre == "宴会場",
		IsAvailable: true,
		IsNew:       ms.Date == "NEW",
		ReleaseDate: pgtype.Date{Time: releaseDate, Valid: true},
		DeleteDate:  song.DeleteDate,
	})
	if updateErr != nil {
		return sqlc.Song{}, fmt.Errorf("failed to update song: %w", updateErr)
	}

	return updatedSong, nil
}

func handleBeatmaps(queries *sqlc.Queries, songID uuid.UUID, ms maimaisong) error {
	// check if std beatmaps exist
	if ms.Basic != "" {
		beatmapType := "std"
		if _, err := upsertBeatmap(queries, songID, "basic", ms.Basic, beatmapType); err != nil {
			return err
		}
		if _, err := upsertBeatmap(queries, songID, "advanced", ms.Advanced, beatmapType); err != nil {
			return err
		}
		if _, err := upsertBeatmap(queries, songID, "expert", ms.Expert, beatmapType); err != nil {
			return err
		}
		if _, err := upsertBeatmap(queries, songID, "master", ms.Master, beatmapType); err != nil {
			return err
		}
		// remaster don't always exist
		if ms.ReMaster != "" {
			if _, err := upsertBeatmap(queries, songID, "remaster", ms.ReMaster, beatmapType); err != nil {
				return err
			}
		}
	}

	// check if dx beatmaps exist
	if ms.DxBasic != "" {
		beatmapType := "dx"
		if _, err := upsertBeatmap(queries, songID, "basic", ms.DxBasic, beatmapType); err != nil {
			return err
		}
		if _, err := upsertBeatmap(queries, songID, "advanced", ms.DxAdvanced, beatmapType); err != nil {
			return err
		}
		if _, err := upsertBeatmap(queries, songID, "expert", ms.DxExpert, beatmapType); err != nil {
			return err
		}
		if _, err := upsertBeatmap(queries, songID, "master", ms.DxMaster, beatmapType); err != nil {
			return err
		}
		// remaster don't always exist
		if ms.DxReMaster != "" {
			if _, err := upsertBeatmap(queries, songID, "remaster", ms.DxReMaster, beatmapType); err != nil {
				return err
			}
		}
	}

	// check if utage beatmaps exist
	if ms.Utage != "" {
		beatmapType := "utage"
		if _, err := upsertBeatmap(queries, songID, "utage", ms.Utage, beatmapType); err != nil {
			return err
		}
	}

	return nil
}

func upsertBeatmap(queries *sqlc.Queries, songID uuid.UUID, difficulty, level, beatmapType string) (sqlc.Beatmap, error) {
	beatmap, getBeatmapErr := queries.GetBeatmapBySongIDDifficultyAndType(context.Background(), sqlc.GetBeatmapBySongIDDifficultyAndTypeParams{
		SongID:     songID,
		Difficulty: difficulty,
		Type:       beatmapType,
	})
	if getBeatmapErr != nil {
		if strings.Contains(getBeatmapErr.Error(), "no rows in result set") {
			// insert if it does not exist in DB
			newBeatmap, createBeatmapErr := queries.CreateBeatmap(context.Background(), sqlc.CreateBeatmapParams{
				BeatmapID:  uuid.New(),
				SongID:     songID,
				Difficulty: difficulty,
				Level:      level,
				Type:       beatmapType,
			})
			if createBeatmapErr != nil {
				return sqlc.Beatmap{}, fmt.Errorf("failed to create beatmap: %w", createBeatmapErr)
			}
			// return newly created song
			return newBeatmap, nil
		}
		return sqlc.Beatmap{}, fmt.Errorf("failed to get beatmap: %w", getBeatmapErr)
	}

	// FIXME: add update

	return beatmap, nil
}

var versionMap = map[string]string{
	"000": "",
	"100": "maimai",
	"110": "maimai PLUS",
	"120": "GreeN",
	"130": "GreeN PLUS",
	"140": "ORANGE",
	"150": "ORANGE PLUS",
	"160": "PiNK",
	"170": "PiNK PLUS",
	"180": "MURASAKi",
	"185": "MURASAKi PLUS",
	"190": "MiLK",
	"195": "MiLK PLUS",
	"199": "FiNALE",
	"200": "maimaiでらっくす",
	"205": "maimaiでらっくす PLUS",
	"210": "Splash",
	"215": "Splash PLUS",
	"220": "UNiVERSE",
	"225": "UNiVERSE PLUS",
	"230": "FESTiVAL",
	"235": "FESTiVAL PLUS",
	"240": "BUDDiES",
	"245": "BUDDiES PLUS",
	"250": "PRiSM",
	"255": "PRiSM PLUS",
}
