package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	database "github.com/asashakira/mai.gg/internal/database/sqlc"
	"github.com/asashakira/mai.gg/internal/utils"
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

	// select * from users
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
