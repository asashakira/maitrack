package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asashakira/mai.gg-api/internal/api/model"
	database "github.com/asashakira/mai.gg-api/internal/database/sqlc"
	"github.com/asashakira/mai.gg-api/internal/scraper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		GameName     string `json:"gameName"`
		TagLine      string `json:"tagLine"`
		SegaID       string `json:"segaID"`
		SegaPassword string `json:"segaPassword"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	// scrape user data to make sure segaID and segaPassword is valid
	u := model.User {
		SegaID: params.SegaID,
		SegaPassword: params.SegaPassword,
	}
	scrapeErr := scraper.FetchUserData(&u)
	if scrapeErr != nil {
		errorMessage := fmt.Sprintf("failed to fetch user data from maimaidxnet: %s", scrapeErr)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	// Hash passwords
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to hash password %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}
	hashedSegaPassword, err := bcrypt.GenerateFromPassword([]byte(params.SegaPassword), bcrypt.DefaultCost)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to hash sega password %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	// insert to users table
	user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
		UserID:       uuid.New(),
		Username:     params.Username,
		Password:     string(hashedPassword),
		SegaID:       params.SegaID,
		SegaPassword: string(hashedSegaPassword),
		GameName:     params.GameName,
		TagLine:      params.TagLine,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("CreateUser %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	// create user data
	userData, err := h.queries.CreateUserData(r.Context(), database.CreateUserDataParams{
		ID:              uuid.New(),
		UserID:          user.UserID,
		GameName:        params.GameName,
		TagLine:         params.TagLine,
		Rating:          u.Rating,
		SeasonPlayCount: u.SeasonPlayCount,
		TotalPlayCount:  u.TotalPlayCount,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error Creating UserData: %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	// create scrape metadata
	defaultLastPlayedAtTime, _ := time.Parse("2006-01-02 15:04", "2006-01-02 15:04")
	_, err = h.queries.CreateUserScrapeMetadata(r.Context(), database.CreateUserScrapeMetadataParams{
		UserID:       user.UserID,
		LastPlayedAt: pgtype.Timestamp{Time: defaultLastPlayedAtTime, Valid: true},
	})
	if err != nil {
		errorMessage := fmt.Sprintf("CreateUserScrapeMetadata %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	respondWithJSON(w, 200, model.ConvertUser(user, userData))
}
