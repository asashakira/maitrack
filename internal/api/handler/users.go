package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	database "github.com/asashakira/maitrack/internal/database/sqlc"
	"github.com/asashakira/maitrack/internal/scraper"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/asashakira/maitrack/pkg/maimaiclient"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) GetUserByMaiID(w http.ResponseWriter, r *http.Request) {
	maiID := chi.URLParam(r, "maiID")

	// get gamename and tagline
	gameName, tagLine, decodeMaiIDErr := decodeMaiID(maiID)
	if decodeMaiIDErr != nil {
		errorMessage := fmt.Sprintf("error decoding maiID %s", decodeMaiIDErr)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	user, err := h.queries.GetUserByMaiID(r.Context(), database.GetUserByMaiIDParams{
		GameName: gameName,
		TagLine:  tagLine,
	})
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No user found with provided fields: %s", err)
			log.Println(errorMessage)
			utils.RespondWithError(w, 404, errorMessage)
			return
		}
		errorMessage := fmt.Sprintf("GetUserByMaiID %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	utils.RespondWithJSON(w, 200, user)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.queries.GetAllUsers(r.Context())
	if err != nil {
		errorMessage := fmt.Sprintf("GetAllUsers %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}
	utils.RespondWithJSON(w, 200, users)
}

func (h *Handler) GetUserHealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, "Hello")
}

func (h *Handler) UpdateUserByMaiID(w http.ResponseWriter, r *http.Request) {
	maiID := chi.URLParam(r, "maiID")

	// get gamename and tagline
	gameName, tagLine, decodeMaiIDErr := decodeMaiID(maiID)
	if decodeMaiIDErr != nil {
		errorMessage := fmt.Sprintf("error decoding maiID %s", decodeMaiIDErr)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	// get user data
	user, err := h.queries.GetUserByMaiID(r.Context(), database.GetUserByMaiIDParams{
		GameName: gameName,
		TagLine:  tagLine,
	})
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No user found with provided fields: %s", err)
			log.Println(errorMessage)
			utils.RespondWithError(w, 404, errorMessage)
			return
		}
		errorMessage := fmt.Sprintf("GetUserByMaiID error: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	segaCreds, err := h.queries.GetSegaCredentials(r.Context(), database.GetSegaCredentialsParams{
		GameName: gameName,
		TagLine:  tagLine,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("GetSegaCredentials error: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	decryptedSegaID, decryptErr := utils.Decrypt(segaCreds.SegaID)
	if decryptErr != nil {
		log.Printf("failed to decrypt SEGA ID: %s", decryptErr)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}
	decryptedSegaPassword, decryptErr := utils.Decrypt(segaCreds.SegaPassword)
	if decryptErr != nil {
		log.Printf("failed to decrypt SEGA password: %s", decryptErr)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}
	m := maimaiclient.New()
	loginErr := m.Login(decryptedSegaID, decryptedSegaPassword)
	if loginErr != nil {
		log.Printf("Failed to login to maimai with SEGAID '%s': %s\n", segaCreds.SegaID, err)
		utils.RespondWithError(w, 400, "Invalid SegaID or password")
		return
	}

	// scrape user and save to database
	scrapeErr := scraper.ScrapeUser(m, h.queries, scraper.ScrapeUserParams{
		UserID:       user.UserID,
		GameName:     user.GameName,
		TagLine:      user.TagLine,
		LastPlayedAt: user.LastPlayedAt,
	})

	if scrapeErr != nil {
		log.Printf("Failed to scrape user data: %s", scrapeErr)
		utils.RespondWithError(w, 400, "Invalid SegaID or password")
		return
	}

	utils.RespondWithJSON(w, 200, "")
}
