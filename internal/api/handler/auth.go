package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/asashakira/mai.gg/internal/scraper"
	"github.com/asashakira/mai.gg/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		SegaID       string `json:"segaID"`
		SegaPassword string `json:"segaPassword"`
		GameName     string `json:"gameName"`
		TagLine      string `json:"tagLine"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	// scrape user data to make sure segaID and segaPassword is valid
	scrapedUserData, scrapeErr := scraper.ScrapeUserData(params.SegaID, params.Password)
	if scrapeErr != nil {
		errorMessage := fmt.Sprintf("failed to scrape user data from maimaidxnet: %s", scrapeErr)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// Hash passwords
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to hash password %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}
	encryptedSegaPassword, err := utils.Encrypt(params.SegaPassword)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to encrypt sega password %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// create user
	user, err := h.queries.CreateUser(r.Context(), database.CreateUserParams{
		UserID:       uuid.New(),
		Username:     params.Username,
		Password:     string(hashedPassword),
		SegaID:       params.SegaID,
		SegaPassword: encryptedSegaPassword,
		GameName:     params.GameName,
		TagLine:      params.TagLine,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error Creating User: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// create user data
	userData, err := h.queries.CreateUserData(r.Context(), database.CreateUserDataParams{
		ID:              uuid.New(),
		UserID:          user.UserID,
		GameName:        user.GameName,
		TagLine:         user.TagLine,
		Rating:          scrapedUserData.Rating,
		SeasonPlayCount: scrapedUserData.SeasonPlayCount,
		TotalPlayCount:  scrapedUserData.TotalPlayCount,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error Creating UserData: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	defaultLastPlayedAtTime, _ := utils.StringToUTCTime("2006-01-02 15:04")
	// create scrape metadata
	usermetadata, err := h.queries.CreateUserMetadata(r.Context(), database.CreateUserMetadataParams{
		UserID:       user.UserID,
		LastPlayedAt: pgtype.Timestamp{Time: defaultLastPlayedAtTime, Valid: true},
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Error Creating UserMetadata: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	data := map[string]interface{}{
		"userID":          user.UserID,
		"username":        user.Username,
		"gameName":        user.GameName,
		"tagLine":         user.TagLine,
		"rating":          userData.Rating,
		"seasonPlayCount": userData.SeasonPlayCount,
		"totalPlayCount":  userData.TotalPlayCount,
		"lastPlayedAt":    usermetadata.LastPlayedAt,
	}

	utils.RespondWithJSON(w, 200, data)
}
