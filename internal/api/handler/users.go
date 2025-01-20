package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/asashakira/mai.gg-api/internal/api/model"
	database "github.com/asashakira/mai.gg-api/internal/database/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) GetUserByMaiID(w http.ResponseWriter, r *http.Request) {
	gameName := chi.URLParam(r, "gameName")
	gameName, err := url.QueryUnescape(gameName)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error decoding gameName from url: %s", err))
		return
	}
	tagLine := chi.URLParam(r, "tagLine")
	tagLine, err = url.QueryUnescape(tagLine)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error decoding tagLine from url: %s", err))
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
			respondWithError(w, 404, errorMessage)
			return
		}
		errorMessage := fmt.Sprintf("GetUserByMaiID %s", err)
		log.Println(errorMessage)
		respondWithError(w, 400, errorMessage)
		return
	}

	// select * from user_data
	userData, err := h.queries.GetUserDataByMaiID(r.Context(), database.GetUserDataByMaiIDParams{
		GameName: chi.URLParam(r, "gameName"),
		TagLine:  chi.URLParam(r, "tagLine"),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("GetUserDataByMaiID Error: %s", err))
		return
	}

	respondWithJSON(w, 200, model.ConvertUser(user, userData))
}

// func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		UserID       *string `json:"userID,omitempty"`
// 		Username     *string `json:"username,omitempty"`
// 		Password     *string `json:"password,omitempty"`
// 		SegaID       *string `json:"segaID,omitempty"`
// 		SegaPassword *string `json:"segaPassword,omitempty"`
// 		GameName     *string `json:"gameName,omitempty"`
// 		TagLine      *string `json:"tagLine,omitempty"`
// 	}
// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %s", err))
// 		return
// 	}
//
// 	// Fetch existing user
// 	userID, parseErr := uuid.Parse(*params.UserID)
// 	if parseErr != nil {
// 		respondWithError(w, 400, fmt.Sprintf("error parsing UserID: %s", err))
// 		return
// 	}
// 	user, err := h.queries.GetUserByID(r.Context(), userID)
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("user not found %s", err)
// 		log.Println(errorMessage)
// 		respondWithError(w, 400, errorMessage)
// 		return
// 	}
//
// 	// Update only the fields provided in the request
// 	updatedUser, err := h.queries.UpdateUser(r.Context(), database.UpdateUserParams{
// 		UserID:       user.UserID,
// 		Username:     ifNotNil(params.Username, user.Username),
// 		Password:     ifNotNil(params.Password, user.Password),
// 		SegaID:       ifNotNil(params.SegaID, user.SegaID),
// 		SegaPassword: ifNotNil(params.SegaPassword, user.SegaPassword),
// 		GameName:     ifNotNil(params.GameName, user.GameName),
// 		TagLine:      ifNotNil(params.TagLine, user.TagLine),
// 	})
// 	if err != nil {
// 		errorMessage := fmt.Sprintf("UpdateUser %s", err)
// 		log.Println(errorMessage)
// 		respondWithError(w, 400, errorMessage)
// 		return
// 	}
//
// 	respondWithJSON(w, 200, model.ConvertUser(updatedUser))
// }
