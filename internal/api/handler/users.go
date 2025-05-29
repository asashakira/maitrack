package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/asashakira/maitrack/internal/scraper"
	"github.com/asashakira/maitrack/internal/utils"
	"github.com/asashakira/maitrack/pkg/maimaiclient"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) GetUserByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	user, err := h.queries.GetUserByUserID(r.Context(), userID)
	if err != nil {
		// Handle "no rows found"
		if errors.Is(err, pgx.ErrNoRows) {
			errorMessage := fmt.Sprintf("No user found with provided fields: %s", err)
			log.Println(errorMessage)
			utils.RespondWithError(w, 404, errorMessage)
			return
		}
		errorMessage := fmt.Sprintf("GetUserByUserID %s", err)
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

func (h *Handler) UpdateUserByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	// get user data
	user, err := h.queries.GetUserByUserID(r.Context(), userID)
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

	segaCreds, err := h.queries.GetSegaCredentialsByUserID(r.Context(), userID)
	if err != nil {
		errorMessage := fmt.Sprintf("GetSegaCredentials error: %s", err)
		log.Println(errorMessage)
		utils.RespondWithError(w, 400, errorMessage)
		return
	}

	decryptedSegaID, decryptErr := utils.Decrypt(segaCreds.EncryptedSegaID)
	if decryptErr != nil {
		log.Printf("failed to decrypt SEGA ID: %s", decryptErr)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}
	decryptedSegaPassword, decryptErr := utils.Decrypt(segaCreds.EncryptedSegaPassword)
	if decryptErr != nil {
		log.Printf("failed to decrypt SEGA password: %s", decryptErr)
		utils.RespondWithError(w, 500, "Internal Server Error")
		return
	}
	m := maimaiclient.New()
	loginErr := m.Login(decryptedSegaID, decryptedSegaPassword)
	if loginErr != nil {
		log.Printf("Failed to login to maimai with SEGAID '%s': %s\n", segaCreds.EncryptedSegaID, err)
		utils.RespondWithError(w, 400, "Invalid SegaID or password")
		return
	}

	// scrape user and save to database
	scrapeErr := scraper.ScrapeUser(m, h.queries, scraper.ScrapeUserParams{
		ID:           user.ID,
		UserID:       user.UserID,
		LastPlayedAt: user.LastPlayedAt,
	})

	if scrapeErr != nil {
		log.Printf("Failed to scrape user data: %s", scrapeErr)
		utils.RespondWithError(w, 400, "Invalid SegaID or password")
		return
	}

	utils.RespondWithJSON(w, 200, "")
}
