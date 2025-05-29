package gamerch

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/asashakira/maitrack/internal/utils"
)

type Song struct {
	ID          string `json:"id,omitempty"`
	AltKey      string `json:"altkey,omitempty"`
	Title       string `json:"title,omitempty"`
	Artist      string `json:"artist,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Bpm         string `json:"bpm,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
	Version     string `json:"version,omitempty"`
	IsUtage     bool   `json:"isUtage,omitempty"`
	IsAvailable bool   `json:"isAvailable,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	DeleteDate  string `json:"deleteDate,omitempty"`
}

func (s *Song) Format() {
	s.Title = RemoveNote(s.Title)
	s.Artist = RemoveNote(s.Artist)
	s.Genre = RemoveNote(s.Genre)
	s.Bpm = RemoveNote(s.Bpm)
	s.ReleaseDate = RemoveNote(s.ReleaseDate)
	s.ReleaseDate = FormatDate(s.ReleaseDate)
	s.DeleteDate = RemoveNote(s.DeleteDate)
	s.DeleteDate = FormatDate(s.DeleteDate)

	s.AltKey = utils.CreateAltKey(s.Title, s.Artist)
}

func PrintSongs(songs []Song) {
	for _, song := range songs {
		PrintSong(song)
	}
}

func PrintSong(s Song) {
	fmt.Printf("%-15s %v\n", "ID:", s.ID)
	fmt.Printf("%-15s %v\n", "AltKey:", s.AltKey)
	fmt.Printf("%-15s %v\n", "Title:", s.Title)
	fmt.Printf("%-15s %v\n", "Artist:", s.Artist)
	fmt.Printf("%-15s %v\n", "Genre:", s.Genre)
	fmt.Printf("%-15s %v\n", "BPM:", s.Bpm)
	fmt.Printf("%-15s %v\n", "ImageUrl:", s.ImageUrl)
	fmt.Printf("%-15s %v\n", "Version:", s.Version)
	fmt.Printf("%-15s %v\n", "IsUtage:", s.IsUtage)
	fmt.Printf("%-15s %v\n", "IsAvailable:", s.IsAvailable)
	fmt.Printf("%-15s %v\n", "ReleaseDate:", s.ReleaseDate)
	fmt.Printf("%-15s %v\n", "DeleteDate:", s.DeleteDate)
	fmt.Println()
}

func DumpSongsAsJson(songs []Song) error {
	jsonByte, err := json.Marshal(songs)
	if err != nil {
		return err
	}

	// write to file
	os.MkdirAll("./tmp/json/", os.ModePerm)
	err = os.WriteFile("./tmp/json/songs.json", jsonByte, 0666)
	if err != nil {
		return err
	}
	return nil
}

type Beatmap struct {
	BeatmapID     string  `json:"beatmapID,omitempty"`
	SongID        string  `json:"songID,omitempty"`
	Difficulty    string  `json:"difficulty,omitempty"`
	Level         string  `json:"level,omitempty"`
	InternalLevel float64 `json:"internalLevel,omitempty"`
	Type          string  `json:"type,omitempty"`
	TotalNotes    int32   `json:"totalNotes,omitempty"`
	Tap           int32   `json:"tap,omitempty"`
	Hold          int32   `json:"hold,omitempty"`
	Slide         int32   `json:"slide,omitempty"`
	Touch         int32   `json:"touch,omitempty"`
	Break         int32   `json:"break,omitempty"`
	NoteDesigner  string  `json:"noteDesigner,omitempty"`
	MaxDxScore    int32   `json:"maxDxScore,omitempty"`
	IsValid       bool    `json:"isValid,omitempty"`
}

func PrintBeatmaps(beatmaps []Beatmap) {
	for _, beatmap := range beatmaps {
		PrintBeatmap(beatmap)
	}
}

func PrintBeatmap(b Beatmap) {
	fmt.Printf("%-15s %v\n", "BeatmapID:", b.BeatmapID)
	fmt.Printf("%-15s %v\n", "SongID:", b.SongID)
	fmt.Printf("%-15s %v\n", "Difficulty:", b.Difficulty)
	fmt.Printf("%-15s %v\n", "Level:", b.Level)
	fmt.Printf("%-15s %v\n", "InternalLevel:", b.InternalLevel)
	fmt.Printf("%-15s %v\n", "Type:", b.Type)
	fmt.Printf("%-15s %v\n", "TotalNotes:", b.TotalNotes)
	fmt.Printf("%-15s %v\n", "Tap:", b.Tap)
	fmt.Printf("%-15s %v\n", "Hold:", b.Hold)
	fmt.Printf("%-15s %v\n", "Slide:", b.Slide)
	fmt.Printf("%-15s %v\n", "Touch:", b.Touch)
	fmt.Printf("%-15s %v\n", "Break:", b.Break)
	fmt.Printf("%-15s %v\n", "NoteDesigner:", b.NoteDesigner)
	fmt.Printf("%-15s %v\n", "MaxDxScore:", b.MaxDxScore)
	fmt.Printf("%-15s %v\n", "IsValid:", b.IsValid)
	fmt.Println()
}

func ParseDifficulty(s string) (string, error) {
	colorToDifficulty := map[string]string{
		"#00ced1": "easy",
		"#98fb98": "basic",
		"#ffa500": "advanced",
		"#fa8080": "expert",
		"#ee82ee": "master",
		"#ffceff": "re:master",
		"#ff5296": "utage",
	}

	re := regexp.MustCompile(`background-color:(#[0-9a-f]+)`)
	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		return "", fmt.Errorf("failed to parse difficulty")
	}
	color := match[1]
	return colorToDifficulty[color], nil
}

func DumpBeatmapsAsJson(beatmaps []Beatmap) error {
	jsonByte, err := json.Marshal(beatmaps)
	if err != nil {
		return err
	}

	// write to file
	os.MkdirAll("./tmp/json/", os.ModePerm)
	err = os.WriteFile("./tmp/json/beatmaps.json", jsonByte, 0666)
	if err != nil {
		return err
	}
	return nil
}
